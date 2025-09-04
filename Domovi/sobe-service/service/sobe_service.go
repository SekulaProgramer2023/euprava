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
	collection := db.Client.Database("euprava").Collection("sobe")

	// Konvertuj string u ObjectID
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, fmt.Errorf("nevalidan roomID: %w", err)
	}

	// Dohvati sobu
	var soba models.Soba
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&soba)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("soba nije pronađena")
		}
		return nil, err
	}

	if !soba.IsFree || soba.Capacity <= 0 {
		return nil, fmt.Errorf("soba je puna")
	}

	soba.Users = append(soba.Users, userID)
	soba.Capacity -= 1
	if soba.Capacity == 0 {
		soba.IsFree = false
	}

	update := bson.M{
		"$set": bson.M{
			"users":    soba.Users,
			"capacity": soba.Capacity,
			"isFree":   soba.IsFree,
		},
	}

	_, err = collection.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		return nil, fmt.Errorf("greška pri update-u sobe: %w", err)
	}

	return &soba, nil
}
