package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"notification-service/db"
	"notification-service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Kreira novu notifikaciju
func CreateNotification(ctx context.Context, notif models.Notification) error {
	collection := db.Client.Database("eupravaM").Collection("notifications")

	_, err := collection.InsertOne(ctx, notif)
	if err != nil {
		return fmt.Errorf("greška pri čuvanju notifikacije: %w", err)
	}

	return nil
}

// Dohvata sve notifikacije
func GetAllNotifications(ctx context.Context) ([]models.Notification, error) {
	collection := db.Client.Database("eupravaM").Collection("notifications")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Notification{}, nil
		}
		return nil, fmt.Errorf("greška pri dohvatanju notifikacija: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.Notification
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("greška pri parsiranju notifikacija: %w", err)
	}

	return results, nil
}

// Dohvata notifikacije po userId
func GetNotificationsByUser(ctx context.Context, userId string) ([]models.Notification, error) {
	collection := db.Client.Database("eupravaM").Collection("notifications")

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Notification{}, nil
		}
		return nil, fmt.Errorf("greška pri dohvatanju notifikacija za usera %s: %w", userId, err)
	}
	defer cursor.Close(ctx)

	var results []models.Notification
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("greška pri parsiranju notifikacija: %w", err)
	}

	return results, nil
}

func CreateJeloRemainingNotification(ctx context.Context, jelovnikID, jelovnikNaziv, jeloID, jeloNaziv string, remaining int, datum time.Time) error {
	if remaining > 2 {
		return nil // nije potrebno slati notifikaciju
	}

	title := fmt.Sprintf("Ostatak jela: %s", jeloNaziv)
	message := fmt.Sprintf("Za jelovnik %s, ostalo je još %d porcije jela %s", jelovnikNaziv, remaining, jeloNaziv)

	notif := models.Notification{
		ID:             primitive.NewObjectID(),
		Title:          title,
		Message:        message,
		Type:           "obrok",
		JelovnikID:     jelovnikID,
		JelovnikNaziv:  jelovnikNaziv,
		JeloID:         jeloID,
		JeloNaziv:      jeloNaziv,
		DatumJelovnika: datum.Format(time.RFC3339), // RFC3339
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	return CreateNotification(ctx, notif)
}
