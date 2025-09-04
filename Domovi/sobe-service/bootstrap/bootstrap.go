package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"sobe-service/db"
	"sobe-service/models"
)

func InsertInitialSobe() {
	if os.Getenv("ENABLE_BOOTSTRAP") != "true" {
		return
	}

	collection := db.Client.Database("euprava").Collection("sobe")

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error counting sobe:", err)
		return
	}

	if count > 0 {
		return // Skip if sobe already exist
	}

	var sobe []interface{}

	// Definiši 5 različitih soba
	soba1 := models.Soba{
		ID:         primitive.NewObjectID(),
		RoomNumber: "101",
		Capacity:   2,
		Users:      []string{},
		OnBudget:   true,
		IsFree:     true,
	}

	soba2 := models.Soba{
		ID:         primitive.NewObjectID(),
		RoomNumber: "102",
		Capacity:   3,
		Users:      []string{},
		OnBudget:   false,
		IsFree:     true,
	}

	soba3 := models.Soba{
		ID:         primitive.NewObjectID(),
		RoomNumber: "201",
		Capacity:   4,
		Users:      []string{},
		OnBudget:   true,
		IsFree:     false,
	}

	soba4 := models.Soba{
		ID:         primitive.NewObjectID(),
		RoomNumber: "202",
		Capacity:   1,
		Users:      []string{},
		OnBudget:   false,
		IsFree:     true,
	}

	soba5 := models.Soba{
		ID:         primitive.NewObjectID(),
		RoomNumber: "301",
		Capacity:   0,
		Users:      []string{},
		OnBudget:   true,
		IsFree:     false,
	}

	sobe = append(sobe, soba1, soba2, soba3, soba4, soba5)

	_, err = collection.InsertMany(context.TODO(), sobe)
	if err != nil {
		fmt.Println("Error inserting initial sobe:", err)
	} else {
		fmt.Println("Inserted 5 initial sobe")
	}
}

func ClearUsers() {

	collection := db.Client.Database("euprava").Collection("sobe")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing users:", err)
	} else {
		fmt.Println("Cleared users from database")
	}
}
