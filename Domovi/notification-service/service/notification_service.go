package services

import (
	"context"
	"fmt"
	"notification-service/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ✅ Kreiranje nove notifikacije
func CreateNotification(ctx context.Context, notif map[string]interface{}) error {
	collection := db.Client.Database("euprava").Collection("notifications")

	_, err := collection.InsertOne(ctx, notif)
	if err != nil {
		return fmt.Errorf("greška pri čuvanju notifikacije: %w", err)
	}

	return nil
}

// ✅ Dohvatanje svih notifikacija
func GetAllNotifications(ctx context.Context) ([]map[string]interface{}, error) {
	collection := db.Client.Database("euprava").Collection("notifications")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("greška pri dohvatanju notifikacija: %w", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("greška pri parsiranju notifikacija: %w", err)
	}

	return results, nil
}

// ✅ Dohvatanje notifikacija po userId
func GetNotificationsByUser(ctx context.Context, userId string) ([]map[string]interface{}, error) {
	collection := db.Client.Database("euprava").Collection("notifications")

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("greška pri dohvatanju notifikacija za usera %s: %w", userId, err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("greška pri parsiranju notifikacija: %w", err)
	}

	return results, nil
}
