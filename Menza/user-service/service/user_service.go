package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"user-service/db"
	"user-service/models"
)

func GetUsers() ([]models.User, error) {
	collection := db.Client.Database("eupravaM").Collection("users")
	var users []models.User
	cursor, err := collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

// RegisterUser registruje korisnika i prosleđuje ga drugom sistemu
func RegisterUser(user models.User) (models.User, error) {
	collection := db.Client.Database("eupravaM").Collection("users")

	// 1. Provera da li korisnik već postoji
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return models.User{}, fmt.Errorf("user with email %s already exists", user.Email)
	} else if err != mongo.ErrNoDocuments {
		return models.User{}, err
	}

	// 2. Sačuvaj originalnu lozinku za forward
	originalPassword := user.Password

	// 3. Hešuj lozinku za MongoDB
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)

	// 4. Sačuvaj u MongoDB
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.User{}, err
	}

	// 5. Vrati originalnu lozinku za forward
	user.Password = originalPassword

	// 6. Marshal u JSON
	body, err := json.Marshal(user)
	if err != nil {
		return models.User{}, err
	}

	// 7. Pošalji POST request drugom sistemu
	resp, err := http.Post("http://host.docker.internal/domovi/users/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return models.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return models.User{}, fmt.Errorf("forward request failed with status: %s", resp.Status)
	}

	return user, nil
}

func LoginUser(user models.User) (models.User, error) {
	collection := db.Client.Database("eupravaM").Collection("users")
	var dbUser models.User
	err := collection.FindOne(context.TODO(), map[string]interface{}{"email": user.Email}).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	fmt.Println(dbUser.Password)
	fmt.Println(user.Password)
	if err != nil {
		return models.User{}, errors.New("invalid password")
	}

	return dbUser, nil
}
