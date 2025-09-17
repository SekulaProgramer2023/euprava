package services

import (
	"bytes"
	"context"
	"dogadjaj-service/db"
	"dogadjaj-service/models"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

	// Pronađemo događaj da dobijemo naziv
	var dogadjaj struct {
		Naziv string `bson:"naziv"`
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&dogadjaj)
	if err != nil {
		return fmt.Errorf("događaj nije pronađen: %w", err)
	}

	// Update statusa
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": status}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("greška pri ažuriranju statusa: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("događaj nije pronađen")
	}

	notification := map[string]interface{}{
		"user_id":     "admin",
		"dogadjaj_id": id,
		"message":     fmt.Sprintf("Status događaja '%s' je promenjen na: %s", dogadjaj.Naziv, status),
		"created_at":  time.Now(),
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("greška pri serijalizaciji notifikacije: %w", err)
	}

	resp, err := http.Post("http://notification-service:8088/notification", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("greška pri slanju notifikacije: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notifikacija nije prihvaćena, status: %d", resp.StatusCode)
	}

	return nil
}
