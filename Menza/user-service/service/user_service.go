package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"user-service/db"
	"user-service/models"
)

var jwtSecret = []byte("1234")

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
func GetUserByID(id string) (models.User, error) {
	collection := db.Client.Database("eupravaM").Collection("users")

	// konvertovanje stringa u ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, err
	}

	// Osiguraj da alergije i omiljena jela nisu nil
	if user.Alergije == nil {
		user.Alergije = []string{}
	}
	if user.OmiljenaJela == nil {
		user.OmiljenaJela = []string{}
	}

	return user, nil
}

func RegisterUser(user models.User) (models.User, error) {
	collection := db.Client.Database("eupravaM").Collection("users")
	fmt.Println("111")

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
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.User{}, err
	}

	// 5. Vrati originalnu lozinku za forward
	user.Password = originalPassword

	// --- KREIRANJE FINANSIJSKE KARTICE PRE "7" ---
	userID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println("1")
	karticaBody, err := json.Marshal(map[string]string{
		"userId":  userID,
		"ime":     user.Name,
		"prezime": user.Surname,
		"index":   "2023/001",
	})
	if err != nil {
		return models.User{}, fmt.Errorf("failed to marshal kartica body: %v", err)
	}
	fmt.Println("2")
	fmt.Println(string(karticaBody))

	// Retry logika za POST
	var karticaResp *http.Response
	var karticaErr error
	for i := 0; i < 5; i++ {
		fmt.Println("Pokušaj kreiranja finansijske kartice...")
		karticaResp, karticaErr = http.Post(
			"http://finansijskakartica-service:8085/kartice",
			"application/json",
			bytes.NewBuffer(karticaBody),
		)
		if karticaErr == nil {
			break
		}
		fmt.Println("Kartica service not ready, retrying in 1s...", karticaErr)
		time.Sleep(1 * time.Second)
	}
	if karticaErr != nil {
		return models.User{}, fmt.Errorf("failed to create kartica after retries: %v", karticaErr)
	}

	defer karticaResp.Body.Close()
	bodyBytes, _ := io.ReadAll(karticaResp.Body)
	fmt.Println("Kartica service response status:", karticaResp.Status)
	fmt.Println("Kartica service response body:", string(bodyBytes))
	if karticaResp.StatusCode != http.StatusOK && karticaResp.StatusCode != http.StatusCreated {
		return models.User{}, fmt.Errorf("kartica request failed with status: %s", karticaResp.Status)
	}

	// 6. Marshal u JSON za domovi-service
	body, err := json.Marshal(user)
	if err != nil {
		return models.User{}, err
	}
	fmt.Println("7")

	// 7. Pošalji POST request domovi-service
	resp, err := http.Post("http://host.docker.internal/domovi/users/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return models.User{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return models.User{}, fmt.Errorf("forward request failed with status: %s", resp.Status)
	}
	fmt.Println("8")

	return user, nil
}

func LoginUser(user models.User) (string, error) {
	collection := db.Client.Database("eupravaM").Collection("users")
	var dbUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("user not found")
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// kreiranje JWT tokena
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": dbUser.ID.Hex(),
		"email":  dbUser.Email,
		"role":   dbUser.Role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // token važi 24h
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func UpdateAlergije(userID string, alergije []string) error {
	collection := db.Client.Database("eupravaM").Collection("users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"alergije": alergije, // menja celu listu
		},
	}

	_, err = collection.UpdateByID(context.TODO(), objID, update)
	return err
}
func UpdateOmiljenaJela(userID string, jela []string) error {
	collection := db.Client.Database("eupravaM").Collection("users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"omiljenaJela": jela, // menja celu listu
		},
	}

	_, err = collection.UpdateByID(context.TODO(), objID, update)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	collection := db.Client.Database("eupravaM").Collection("users")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}
