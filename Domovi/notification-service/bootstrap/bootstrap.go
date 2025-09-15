package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"notification-service/db"
)

func ClearNotifications() {

	collection := db.Client.Database("euprava").Collection("notifications")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing reviews:", err)
	} else {
		fmt.Println("Cleared review from database")
	}
}
