package services

import (
	"context"
	"dogadjaj-service/db"
	"dogadjaj-service/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Kreiranje događaja
func CreateDogadjaj(dogadjaj models.Dogadjaj) (*models.Dogadjaj, error) {
	collection := db.Client.Database("euprava").Collection("dogadjaj")

	// Ako nije setovan datum slanja zahteva, dodaćemo trenutni
	if dogadjaj.DatumSlanjaZahteva.IsZero() {
		dogadjaj.DatumSlanjaZahteva = time.Now()
	}
	if dogadjaj.Status == "" {
		dogadjaj.Status = "na čekanju"
	}

	_, err := collection.InsertOne(context.TODO(), dogadjaj)
	if err != nil {
		return nil, fmt.Errorf("greška pri kreiranju događaja: %w", err)
	}

	return &dogadjaj, nil
}

// Dohvatanje svih događaja
func GetAllDogadjaji() ([]models.Dogadjaj, error) {
	collection := db.Client.Database("euprava").Collection("dogadjaj")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("greška pri dohvatanju događaja: %w", err)
	}
	defer cursor.Close(context.TODO())

	var dogadjaji []models.Dogadjaj
	if err := cursor.All(context.TODO(), &dogadjaji); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju događaja: %w", err)
	}

	return dogadjaji, nil
}

// Ažuriranje statusa događaja
func UpdateDogadjajStatus(id string, status string) error {
	collection := db.Client.Database("euprava").Collection("dogadjaj")

	// Dozvoljeni statusi
	if status != "prihvaćen" && status != "odbijen" {
		return fmt.Errorf("nevalidan status: %s", status)
	}

	// Pretvaramo string id u ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("nevalidan ID: %w", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("greška pri ažuriranju statusa: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("događaj nije pronađen")
	}

	return nil
}
