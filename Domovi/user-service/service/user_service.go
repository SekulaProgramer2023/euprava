package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"user-service/db"
	"user-service/models"
)

func GetUsers() ([]models.User, error) {
	collection := db.Client.Database("euprava").Collection("users")
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

func RegisterUser(user models.User) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)

	collection := db.Client.Database("euprava").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func LoginUser(user models.User) (models.User, error) {
	collection := db.Client.Database("euprava").Collection("users")
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
