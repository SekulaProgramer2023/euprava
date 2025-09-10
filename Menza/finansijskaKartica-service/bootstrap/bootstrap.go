package bootstrap

import (
	"context"
	"finansijskaKartica-service/db"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func ClearKartice() {
	collection := db.Client.Database("eupravaM").Collection("finansijske_kartice")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing users:", err)
	} else {
		fmt.Println("Cleared users from database")
	}
}
