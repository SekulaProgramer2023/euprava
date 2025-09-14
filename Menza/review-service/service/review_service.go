package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"review-service/db"
	"review-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// helper da proveri da li je ObjectID validan
func isValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// Provera da li postoji korisnik
func userExists(userID string) (bool, error) {
	if !isValidObjectID(userID) {
		return false, fmt.Errorf("nevalidan userID")
	}

	resp, err := http.Get("http://user-service:8081/users")
	if err != nil {
		return false, fmt.Errorf("greška pri pozivu user-service: %w", err)
	}
	defer resp.Body.Close()

	var users []models.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return false, fmt.Errorf("greška pri parsiranju usera: %w", err)
	}

	for _, u := range users {
		if u.ID == userID {
			return true, nil
		}
	}

	return false, nil
}

// Provera da li postoji soba
func jeloExists(jeloID string) (bool, error) {
	if !isValidObjectID(jeloID) {
		return false, fmt.Errorf("nevalidan jeloID")
	}

	resp, err := http.Get("http://jelovnik-service:8083/jela")
	if err != nil {
		return false, fmt.Errorf("greška pri pozivu jelo-service: %w", err)
	}
	defer resp.Body.Close()

	var jelo []models.Jelo
	if err := json.NewDecoder(resp.Body).Decode(&jelo); err != nil {
		return false, fmt.Errorf("greška pri parsiranju jela: %w", err)
	}

	for _, j := range jelo {
		if j.ID == jeloID {
			return true, nil
		}
	}

	return false, nil
}

// Kreiranje review-a sa validacijom
func CreateReview(review models.Review) (*models.Review, error) {
	// provera usera
	exists, err := userExists(review.UserId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("korisnik sa ID %s ne postoji", review.UserId)
	}

	// provera sobe
	exists, err = jeloExists(review.JeloId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("jelo sa ID %s ne postoji", review.JeloId)
	}

	// definisanje kolekcije unutar funkcije
	collection := db.Client.Database("eupravaM").Collection("reviews")

	// insert u MongoDB
	_, err = collection.InsertOne(context.TODO(), review)
	if err != nil {
		return nil, fmt.Errorf("greška pri čuvanju review-a: %w", err)
	}

	return &review, nil
}

// ✅ Dohvatanje svih review-a
func GetAllReviews() ([]models.Review, error) {
	collection := db.Client.Database("eupravaM").Collection("reviews")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("greška pri dohvatanju review-a: %w", err)
	}
	defer cursor.Close(context.TODO())

	var reviews []models.Review
	if err := cursor.All(context.TODO(), &reviews); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju review-a: %w", err)
	}

	return reviews, nil
}

// ✅ Dohvatanje review-a po sobi
func GetReviewsByJeloID(jeloID string) ([]models.Review, error) {
	collection := db.Client.Database("eupravaM").Collection("reviews")

	cursor, err := collection.Find(context.TODO(), bson.M{"jeloId": jeloID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Review{}, nil
		}
		return nil, fmt.Errorf("greška pri dohvatanju review-a za sobu: %w", err)
	}
	defer cursor.Close(context.TODO())

	var reviews []models.Review
	if err := cursor.All(context.TODO(), &reviews); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju review-a: %w", err)
	}

	return reviews, nil
}

func CalculateAverageRating(jeloID string) (float64, error) {
	collection := db.Client.Database("eupravaM").Collection("reviews")

	cursor, err := collection.Find(context.TODO(), bson.M{"jeloId": jeloID})
	if err != nil {
		return 0, fmt.Errorf("greška pri dohvatanju review-a: %w", err)
	}
	defer cursor.Close(context.TODO())

	var reviews []models.Review
	if err := cursor.All(context.TODO(), &reviews); err != nil {
		return 0, fmt.Errorf("greška pri dekodiranju review-a: %w", err)
	}

	if len(reviews) == 0 {
		return 0, nil // nema review-a za sobu
	}

	sum := 0
	for _, r := range reviews {
		sum += r.Rating
	}

	average := float64(sum) / float64(len(reviews))
	return average, nil
}
