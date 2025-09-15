package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"review-service/db"
	"review-service/models"
	"time"

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

	resp, err := http.Get("http://user-service:8080/users")
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
func sobaExists(sobaID string) (bool, error) {
	if !isValidObjectID(sobaID) {
		return false, fmt.Errorf("nevalidan sobaID")
	}

	resp, err := http.Get("http://sobe-service:8082/sobe")
	if err != nil {
		return false, fmt.Errorf("greška pri pozivu sobe-service: %w", err)
	}
	defer resp.Body.Close()

	var sobe []models.Soba
	if err := json.NewDecoder(resp.Body).Decode(&sobe); err != nil {
		return false, fmt.Errorf("greška pri parsiranju soba: %w", err)
	}

	for _, s := range sobe {
		if s.ID == sobaID {
			return true, nil
		}
	}

	return false, nil
}

// CreateReview kreira review i šalje notifikaciju preko REST API-ja
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
	exists, err = sobaExists(review.SobaId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("soba sa ID %s ne postoji", review.SobaId)
	}

	// definisanje kolekcije za review-e
	collection := db.Client.Database("euprava").Collection("reviews")

	// insert u MongoDB
	_, err = collection.InsertOne(context.TODO(), review)
	if err != nil {
		return nil, fmt.Errorf("greška pri čuvanju review-a: %w", err)
	}

	// --- priprema notifikacije sa imenom i prezimenom korisnika i brojem sobe ---

	// Dohvati informacije o korisniku
	var userFullName string
	respUser, err := http.Get(fmt.Sprintf("http://user-service:8080/users/%s", review.UserId))
	if err != nil {
		return &review, fmt.Errorf("greška pri dohvatanju korisnika: %w", err)
	}
	defer respUser.Body.Close()

	if respUser.StatusCode == http.StatusOK {
		var user models.User
		if err := json.NewDecoder(respUser.Body).Decode(&user); err == nil {
			userFullName = fmt.Sprintf("%s %s", user.Name, user.Surname)
		} else {
			userFullName = "Nepoznati korisnik"
		}
	} else {
		userFullName = "Nepoznati korisnik"
	}

	// Dohvati informacije o sobi
	var roomNumber string
	respRooms, err := http.Get("http://sobe-service:8082/sobe")
	if err != nil {
		return &review, fmt.Errorf("greška pri dohvatanju soba: %w", err)
	}
	defer respRooms.Body.Close()

	if respRooms.StatusCode == http.StatusOK {
		var sobe []models.Soba
		if err := json.NewDecoder(respRooms.Body).Decode(&sobe); err == nil {
			// traži sobu sa odgovarajućim ID-om
			found := false
			for _, s := range sobe {
				if s.ID == review.SobaId {
					fmt.Println(s)
					roomNumber = s.RoomNumber
					fmt.Println(roomNumber)
					found = true
					break
				}
			}
			if !found {
				roomNumber = review.SobaId
			}
		} else {
			roomNumber = review.SobaId
		}
	} else {
		roomNumber = review.SobaId
	}

	// Kreiraj notifikaciju sa imenom i brojem sobe
	notification := map[string]interface{}{
		"user_id":    review.UserId,
		"soba_id":    review.SobaId,
		"message":    fmt.Sprintf("Korisnik %s je ostavio recenziju za sobu %s", userFullName, roomNumber),
		"rating":     review.Rating,
		"created_at": time.Now(),
	}

	// serijalizacija u JSON
	data, err := json.Marshal(notification)
	if err != nil {
		return &review, fmt.Errorf("greška pri serijalizaciji notifikacije: %w", err)
	}

	// POST zahtev ka notification-service
	resp, err := http.Post("http://notification-service:8088/notification", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return &review, fmt.Errorf("greška pri slanju notifikacije: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return &review, fmt.Errorf("notification-service je vratio status: %d", resp.StatusCode)
	}

	return &review, nil
}

// ✅ Dohvatanje svih review-a
func GetAllReviews() ([]models.Review, error) {
	collection := db.Client.Database("euprava").Collection("reviews")

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
func GetReviewsBySobaID(sobaID string) ([]models.Review, error) {
	collection := db.Client.Database("euprava").Collection("reviews")

	cursor, err := collection.Find(context.TODO(), bson.M{"soba_id": sobaID})
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

func CalculateAverageRating(sobaID string) (float64, error) {
	collection := db.Client.Database("euprava").Collection("reviews")

	cursor, err := collection.Find(context.TODO(), bson.M{"soba_id": sobaID})
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
