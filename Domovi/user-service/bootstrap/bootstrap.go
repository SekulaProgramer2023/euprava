package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"os"
	"user-service/db"
	"user-service/models"
)

func InsertInitialUsers() {
	if os.Getenv("ENABLE_BOOTSTRAP") != "true" {
		return
	}

	collection := db.Client.Database("euprava").Collection("users")

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error counting users:", err)
		return
	}

	if count > 0 {
		return // Skip if users already exist
	}

	// Dodaj unapred definisane korisnike
	var users []interface{}

	// Dodavanje korisnika "aca"
	hashedPasswordAca, err := bcrypt.GenerateFromPassword([]byte("Aca2024!"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password for aca:", err)
		return
	}
	acaUser := models.User{
		Password: string(hashedPasswordAca),
		Role:     "Manager",
		Name:     "Aca1",
		Surname:  "Admin",
		Email:    "aca@example.com",
		IsActive: true,
	}
	users = append(users, acaUser)

	_, err = collection.InsertMany(context.TODO(), users)
	if err != nil {
		fmt.Println("Error inserting initial users:", err)
	} else {
		fmt.Println("Inserted initial users including 'aca' and 'ana'")
	}
}

func ClearUsers() {

	collection := db.Client.Database("euprava").Collection("users")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing users:", err)
	} else {
		fmt.Println("Cleared users from database")
	}
}
