package bootstrap

import (
	"context"
	"dogadjaj-service/db"
	"dogadjaj-service/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"time"
)

func InsertInitialDogadjaji() {
	if os.Getenv("ENABLE_BOOTSTRAP") != "true" {
		return
	}

	collection := db.Client.Database("euprava").Collection("dogadjaj")

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error counting dogadjaji:", err)
		return
	}

	if count > 0 {
		return // već postoje, ne insertuje ponovo
	}

	var dogadjaji []interface{}

	d1 := models.NewDogadjaj(
		"Studentska žurka",
		"Žurka za sve studente u dvorištu doma.",
		time.Date(2025, 9, 20, 20, 0, 0, 0, time.UTC),
		"meso",
	)

	d2 := models.NewDogadjaj(
		"Sportski turnir",
		"Turnir u malom fudbalu između domova.",
		time.Date(2025, 10, 5, 15, 0, 0, 0, time.UTC), "vege",
	)

	d3 := models.NewDogadjaj(
		"Kviz veče",
		"Opšte znanje, timovi po 4 člana.",
		time.Date(2025, 11, 12, 19, 0, 0, 0, time.UTC),
		"kobasice")

	dogadjaji = append(dogadjaji, d1, d2, d3)

	_, err = collection.InsertMany(context.TODO(), dogadjaji)
	if err != nil {
		fmt.Println("Error inserting initial dogadjaji:", err)
	} else {
		fmt.Println("Inserted 3 initial dogadjaji")
	}
}

func ClearDogadjaj() {
	collection := db.Client.Database("euprava").Collection("dogadjaj")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing dogadjaji:", err)
	} else {
		fmt.Println("Cleared dogadjaji from database")
	}
}
