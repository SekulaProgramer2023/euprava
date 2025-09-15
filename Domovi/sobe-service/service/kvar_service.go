package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"sobe-service/db"
	"sobe-service/models"
	"time"
)

func userExists(userID string) (bool, error) {
	collection := db.Client.Database("euprava").Collection("users")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("nevalidan userID: %w", err)
	}

	count, err := collection.CountDocuments(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ✅ Provera da li postoji soba
func sobaExists(sobaID string) (bool, error) {
	collection := db.Client.Database("euprava").Collection("sobe")
	objID, err := primitive.ObjectIDFromHex(sobaID)
	if err != nil {
		return false, fmt.Errorf("nevalidan sobaID: %w", err)
	}

	count, err := collection.CountDocuments(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateKvar(kvar models.Kvar) (*models.Kvar, error) {
	log.Println("Start CreateKvar")

	// Provera korisnika
	exists, err := userExists(kvar.UserId)
	if err != nil {
		log.Printf("Error checking user: %v\n", err)
		return nil, err
	}
	if !exists {
		errMsg := fmt.Sprintf("korisnik sa ID %s ne postoji", kvar.UserId)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Provera sobe
	exists, err = sobaExists(kvar.SobaId)
	if err != nil {
		log.Printf("Error checking soba: %v\n", err)
		return nil, err
	}
	if !exists {
		errMsg := fmt.Sprintf("soba sa ID %s ne postoji", kvar.SobaId)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Insert u MongoDB
	collection := db.Client.Database("euprava").Collection("kvarovi")
	_, err = collection.InsertOne(context.TODO(), kvar)
	if err != nil {
		log.Printf("Error inserting kvar: %v\n", err)
		return nil, fmt.Errorf("greška pri kreiranju kvara: %w", err)
	}
	log.Println("Kvar successfully inserted in DB")

	// Dohvati ime i prezime korisnika
	userFullName := "Nepoznati korisnik"
	respUser, err := http.Get(fmt.Sprintf("http://user-service:8080/users/%s", kvar.UserId))
	if err == nil && respUser.StatusCode == http.StatusOK {
		defer respUser.Body.Close()
		var user models.User
		if err := json.NewDecoder(respUser.Body).Decode(&user); err == nil {
			userFullName = fmt.Sprintf("%s %s", user.Name, user.Surname)
		} else {
			log.Printf("Error decoding user: %v\n", err)
		}
	} else {
		log.Printf("Error fetching user or bad status: %v, status: %v\n", err, respUser)
	}

	// Dohvati broj sobe
	roomNumber := kvar.SobaId
	respRooms, err := http.Get("http://sobe-service:8082/sobe")
	if err == nil && respRooms.StatusCode == http.StatusOK {
		defer respRooms.Body.Close()
		var sobe []models.Soba
		if err := json.NewDecoder(respRooms.Body).Decode(&sobe); err == nil {
			sobaObjID, _ := primitive.ObjectIDFromHex(kvar.SobaId)
			for _, s := range sobe {
				if s.ID == sobaObjID {
					roomNumber = s.RoomNumber
					break
				}
			}
		} else {
			log.Printf("Error decoding sobe: %v\n", err)
		}
	} else {
		log.Printf("Error fetching sobe or bad status: %v, status: %v\n", err, respRooms)
	}

	// Kreiraj notifikaciju
	notification := map[string]interface{}{
		"user_id":    kvar.UserId,
		"soba_id":    kvar.SobaId,
		"message":    fmt.Sprintf("Korisnik %s je prijavio kvar za sobu %s", userFullName, roomNumber),
		"created_at": time.Now(),
	}

	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling notification: %v\n", err)
		return &kvar, fmt.Errorf("greška pri serijalizaciji notifikacije: %w", err)
	}

	respNotif, err := http.Post("http://notification-service:8088/notification", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending notification: %v\n", err)
		return &kvar, fmt.Errorf("greška pri slanju notifikacije: %w", err)
	}
	defer respNotif.Body.Close()

	if respNotif.StatusCode != http.StatusCreated {
		log.Printf("Notification service returned status: %d\n", respNotif.StatusCode)
		return &kvar, fmt.Errorf("notification-service je vratio status: %d", respNotif.StatusCode)
	}

	log.Println("Notification successfully sent")
	return &kvar, nil
}

// ✅ Dohvatanje svih kvarova
func GetAllKvarovi() ([]models.Kvar, error) {
	collection := db.Client.Database("euprava").Collection("kvarovi")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("greška pri dohvatanju kvarova: %w", err)
	}
	defer cursor.Close(context.TODO())

	var kvarovi []models.Kvar
	if err := cursor.All(context.TODO(), &kvarovi); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju kvarova: %w", err)
	}

	return kvarovi, nil
}

// ✅ Dohvatanje kvarova po sobi
func GetKvaroviBySobaID(sobaID string) ([]models.Kvar, error) {
	collection := db.Client.Database("euprava").Collection("kvarovi")

	// Validacija ObjectID-a (ako koristiš hex ID)
	_, err := primitive.ObjectIDFromHex(sobaID)
	if err != nil {
		return nil, fmt.Errorf("nevalidan sobaID: %w", err)
	}

	cursor, err := collection.Find(context.TODO(), bson.M{"soba_id": sobaID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Kvar{}, nil // nema kvarova
		}
		return nil, fmt.Errorf("greška pri dohvatanju kvarova za sobu: %w", err)
	}
	defer cursor.Close(context.TODO())

	var kvarovi []models.Kvar
	if err := cursor.All(context.TODO(), &kvarovi); err != nil {
		return nil, fmt.Errorf("greška pri dekodiranju kvarova: %w", err)
	}

	return kvarovi, nil
}

func ResolveKvar(kvarID string) error {
	collection := db.Client.Database("euprava").Collection("kvarovi")

	// Validacija ID-a
	objID, err := primitive.ObjectIDFromHex(kvarID)
	if err != nil {
		return fmt.Errorf("nevalidan kvarID: %w", err)
	}

	// Dohvati kvar pre update-a (da imamo info za notifikaciju)
	var kvar models.Kvar
	if err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&kvar); err != nil {
		return fmt.Errorf("kvar sa ID %s nije pronađen: %w", kvarID, err)
	}

	// Update dokumenta (status = true)
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		return fmt.Errorf("greška pri ažuriranju kvara: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("kvar sa ID %s nije pronađen", kvarID)
	}

	// --- Kreiranje notifikacije ---
	// Dohvati ime i prezime korisnika
	var userFullName string
	respUser, err := http.Get(fmt.Sprintf("http://user-service:8080/users/%s", kvar.UserId))
	if err == nil {
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
	} else {
		userFullName = "Nepoznati korisnik"
	}

	// Dohvati broj sobe (sličan princip kao ranije)
	var roomNumber string
	respRoom, err := http.Get("http://sobe-service:8082/sobe")
	if err == nil {
		defer respRoom.Body.Close()
		var sobe []models.Soba
		if err := json.NewDecoder(respRoom.Body).Decode(&sobe); err == nil {
			for _, s := range sobe {
				if s.ID.Hex() == kvar.SobaId {
					roomNumber = s.RoomNumber
					break
				}
			}
		}
	}

	notification := map[string]interface{}{
		"user_id":    kvar.UserId,
		"soba_id":    kvar.SobaId,
		"message":    fmt.Sprintf("Kvar u sobi %s prijavljen od strane %s je otklonjen", roomNumber, userFullName),
		"created_at": time.Now(),
	}

	data, err := json.Marshal(notification)
	if err == nil {
		// Pošalji notifikaciju ka notification-service
		resp, err := http.Post("http://notification-service:8088/notification", "application/json", bytes.NewBuffer(data))
		if err == nil {
			defer resp.Body.Close()
		}
	}

	return nil
}
