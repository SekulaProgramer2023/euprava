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

func userExists(userID string) (bool, error) {
	collection := db.Client.Database("euprava").Collection("users")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("nevalidan userID: %w", err)
	}

	count, err := collection.CountDocuments(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ✅ Provera da li postoji soba
func sobaExists(sobaID string) (bool, error) {
	collection := db.Client.Database("euprava").Collection("sobe")
	objID, err := primitive.ObjectIDFromHex(sobaID)
	if err != nil {
		return false, fmt.Errorf("nevalidan sobaID: %w", err)
	}

	count, err := collection.CountDocuments(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ✅ Kreiranje novog kvara (sa proverom usera i sobe)
func CreateKvar(kvar models.Kvar) (*models.Kvar, error) {
	// Provera korisnika
	exists, err := userExists(kvar.UserId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("korisnik sa ID %s ne postoji", kvar.UserId)
	}

	// Provera sobe
	exists, err = sobaExists(kvar.SobaId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("soba sa ID %s ne postoji", kvar.SobaId)
	}

	collection := db.Client.Database("euprava").Collection("kvarovi")

	_, err = collection.InsertOne(context.TODO(), kvar)
	if err != nil {
		return nil, fmt.Errorf("greška pri kreiranju kvara: %w", err)
	}

	return &kvar, nil
}

// ✅ Dohvatanje svih kvarova
func GetAllKvarovi() ([]models.Kvar, error) {
	collection := db.Client.Database("euprava").Collection("kvarovi")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("greška pri dohvatanju kvarova: %w", err)
	}
	defer cursor.Close(context.TODO())

	var kvarovi []models.Kvar
	if err := cursor.All(context.TODO(), &kvarovi); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju kvarova: %w", err)
	}

	return kvarovi, nil
}

// ✅ Dohvatanje kvarova po sobi
func GetKvaroviBySobaID(sobaID string) ([]models.Kvar, error) {
	collection := db.Client.Database("euprava").Collection("kvarovi")

	// Validacija ObjectID-a (ako koristiš hex ID)
	_, err := primitive.ObjectIDFromHex(sobaID)
	if err != nil {
		return nil, fmt.Errorf("nevalidan sobaID: %w", err)
	}

	cursor, err := collection.Find(context.TODO(), bson.M{"soba_id": sobaID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Kvar{}, nil // nema kvarova
		}
		return nil, fmt.Errorf("greška pri dohvatanju kvarova za sobu: %w", err)
	}
	defer cursor.Close(context.TODO())

	var kvarovi []models.Kvar
	if err := cursor.All(context.TODO(), &kvarovi); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju kvarova: %w", err)
	}

	return kvarovi, nil
}
