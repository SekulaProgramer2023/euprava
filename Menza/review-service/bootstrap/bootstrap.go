package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"review-service/db"
)

func ClearReviews() {

	collection := db.Client.Database("eupravaM").Collection("reviews")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing reviews:", err)
	} else {
		fmt.Println("Cleared review from database")
	}
}
