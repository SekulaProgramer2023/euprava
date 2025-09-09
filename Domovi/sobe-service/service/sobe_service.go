package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sobe-service/db"
	"sobe-service/models"
)

func GetAllSobe() ([]models.Soba, error) {
	collection := db.Client.Database("euprava").Collection("sobe")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("error fetching sobe: %w", err)
	}
	defer cursor.Close(context.TODO())

	var sobe []models.Soba
	if err := cursor.All(context.TODO(), &sobe); err != nil {
		return nil, fmt.Errorf("error decoding sobe: %w", err)
	}

	return sobe, nil
}

// Kreiranje nove sobe
func CreateSoba(soba models.Soba) (*models.Soba, error) {
	collection := db.Client.Database("euprava").Collection("sobe")
	_, err := collection.InsertOne(context.TODO(), soba)
	if err != nil {
		return nil, fmt.Errorf("error inserting soba: %w", err)
	}
	return &soba, nil
}

// Dohvatanje svih soba sa kapacitetom > 0
func GetSobeWithCapacity() ([]models.Soba, error) {
	collection := db.Client.Database("euprava").Collection("sobe")
	filter := bson.M{"capacity": bson.M{"$gt": 0}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching sobe: %w", err)
	}
	defer cursor.Close(context.TODO())

	var sobe []models.Soba
	if err := cursor.All(context.TODO(), &sobe); err != nil {
		return nil, fmt.Errorf("error decoding sobe: %w", err)
	}

	return sobe, nil
}

func UseliUsera(roomID string, userID string) (*models.Soba, error) {
	sobeCollection := db.Client.Database("euprava").Collection("sobe")
	usersCollection := db.Client.Database("euprava").Collection("users")

	// Proveri da li je user već useljen u neku sobu
	var existingSoba models.Soba
	err := sobeCollection.FindOne(context.TODO(), bson.M{"users": userID}).Decode(&existingSoba)
	if err == nil {
		return nil, fmt.Errorf("Korisnik je već useljen u sobu broj %s", existingSoba.RoomNumber)
	} else if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("greška pri proveri korisnika: %w", err)
	}

	// Konvertuj roomID u ObjectID
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, fmt.Errorf("nevalidan roomID: %w", err)
	}

	// Dohvati sobu
	var soba models.Soba
	err = sobeCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&soba)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("soba nije pronađena")
		}
		return nil, err
	}

	// Provera kapaciteta
	if soba.Capacity <= 0 {
		soba.IsFree = false // osiguravamo da je soba "full"
		_, _ = sobeCollection.UpdateByID(context.TODO(), objID, bson.M{"$set": bson.M{"isFree": false}})
		return nil, fmt.Errorf("soba je puna")
	}

	// Useli korisnika u sobu
	soba.Users = append(soba.Users, userID)
	soba.Capacity -= 1

	// Ažuriranje IsFree nakon useljavanja
	if soba.Capacity == 0 {
		soba.IsFree = false
	}

	// Update sobe u bazi
	_, err = sobeCollection.UpdateByID(context.TODO(), objID, bson.M{
		"$set": bson.M{
			"users":    soba.Users,
			"capacity": soba.Capacity,
			"IsFree":   soba.IsFree,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("greška pri update-u sobe: %w", err)
	}

	// Upisi sobu kod usera
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("nevalidan userID: %w", err)
	}

	_, err = usersCollection.UpdateByID(context.TODO(), userObjID, bson.M{
		"$set": bson.M{"soba": objID},
	})
	if err != nil {
		return nil, fmt.Errorf("greška pri update-u korisnika: %w", err)
	}

	return &soba, nil
}

func GetSobaByID(roomID string) (*models.Soba, error) {
	collection := db.Client.Database("euprava").Collection("sobe")

	// Konvertuj string u ObjectID
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, fmt.Errorf("nevalidan roomID: %w", err)
	}

	// Pronađi sobu
	var soba models.Soba
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&soba)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("soba nije pronađena")
		}
		return nil, fmt.Errorf("greška pri dohvaćanju sobe: %w", err)
	}

	return &soba, nil
}
