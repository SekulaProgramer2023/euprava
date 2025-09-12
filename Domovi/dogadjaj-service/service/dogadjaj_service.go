package services

import (
	"context"
	"dogadjaj-service/db"
	"dogadjaj-service/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
