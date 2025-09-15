package service

import (
	"bytes"
	"context"
	"encoding/json"
	"finansijskaKartica-service/models"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FinansijskaKarticaService struct {
	Collection *mongo.Collection
}

// Konstruktor: prima bazu koja je već inicijalizovana
func NewFinansijskaKarticaService(db *mongo.Database) *FinansijskaKarticaService {
	return &FinansijskaKarticaService{
		Collection: db.Collection("finansijske_kartice"),
	}
}

// Kreira novu karticu
func (s *FinansijskaKarticaService) CreateKartica(kartica models.FinansijskaKartica) (models.FinansijskaKartica, error) {
	_, err := s.Collection.InsertOne(context.TODO(), kartica)
	if err != nil {
		return models.FinansijskaKartica{}, err
	}
	return kartica, nil
}

// Vraća sve kartice
func (s *FinansijskaKarticaService) GetKartice() ([]models.FinansijskaKartica, error) {
	var kartice []models.FinansijskaKartica
	cursor, err := s.Collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &kartice); err != nil {
		return nil, err
	}
	return kartice, nil
}

// Uplata novca na karticu
func (s *FinansijskaKarticaService) Deposit(userID primitive.ObjectID, novac float64) (models.FinansijskaKartica, error) {
	filter := bson.M{"userId": userID}
	update := bson.M{"$inc": bson.M{"novac": novac}}

	var updated models.FinansijskaKartica
	err := s.Collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After), // << ovo je ključno
	).Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return models.FinansijskaKartica{}, fmt.Errorf("kartica not found for user %s", userID.Hex())
	}
	if err != nil {
		return models.FinansijskaKartica{}, err
	}
	return updated, nil
}

// NOVO: Vraća karticu po userId
func (s *FinansijskaKarticaService) GetKarticaByUserID(userID primitive.ObjectID) (models.FinansijskaKartica, error) {
	var kartica models.FinansijskaKartica
	err := s.Collection.FindOne(context.TODO(), bson.M{"userId": userID}).Decode(&kartica)
	if err == mongo.ErrNoDocuments {
		return models.FinansijskaKartica{}, fmt.Errorf("kartica for user %s not found", userID.Hex())
	} else if err != nil {
		return models.FinansijskaKartica{}, err
	}
	return kartica, nil
}

// Kupovina doručka (70 RSD)
func (s *FinansijskaKarticaService) buyMeals(userID primitive.ObjectID, cena float64, field string, count int) (models.FinansijskaKartica, error) {
	var kartica models.FinansijskaKartica
	err := s.Collection.FindOne(context.TODO(), bson.M{"userId": userID}).Decode(&kartica)
	if err == mongo.ErrNoDocuments {
		return models.FinansijskaKartica{}, fmt.Errorf("kartica not found for user %s", userID.Hex())
	}
	if err != nil {
		return models.FinansijskaKartica{}, err
	}

	ukupnaCena := cena * float64(count)
	if kartica.Novac < ukupnaCena {
		return models.FinansijskaKartica{}, fmt.Errorf("nedovoljno sredstava (balans: %.2f RSD, potrebno: %.2f RSD)", kartica.Novac, ukupnaCena)
	}

	update := bson.M{
		"$inc": bson.M{
			"novac": -ukupnaCena,
			field:   count,
		},
	}

	var updated models.FinansijskaKartica
	err = s.Collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"userId": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After), // << ovo je ključno
	).Decode(&updated)

	if err != nil {
		return models.FinansijskaKartica{}, err
	}

	return updated, nil
}

func (s *FinansijskaKarticaService) BuyRuckovi(userID primitive.ObjectID, count int) (models.FinansijskaKartica, error) {
	return s.buyMeals(userID, 120, "rucakCount", count)
}

func (s *FinansijskaKarticaService) BuyVecere(userID primitive.ObjectID, count int) (models.FinansijskaKartica, error) {
	return s.buyMeals(userID, 90, "veceraCount", count)
}

func (s *FinansijskaKarticaService) BuyDorucak(userID primitive.ObjectID, count int) (models.FinansijskaKartica, error) {
	return s.buyMeals(userID, 70, "dorucakCount", count)
}
func (s *FinansijskaKarticaService) IskoristiObrok(userID, jelovnikID, jeloID string) (models.FinansijskaKartica, error) {
	// 1. Dohvati sve jelovnike preko REST poziva
	jelovnici, err := GetJelovnici()
	if err != nil {
		return models.FinansijskaKartica{}, err
	}

	// 2. Pronađi jelovnik po ID-u
	var jelovnikPravi *models.Jelovnik
	for _, j := range jelovnici {
		if j.JelovnikID == jelovnikID {
			jelovnikPravi = &j
			break
		}
	}
	if jelovnikPravi == nil {
		return models.FinansijskaKartica{}, fmt.Errorf("jelovnik sa ID %s nije pronađen", jelovnikID)
	}

	// 3. Pronađi jelo po ID-u unutar tog jelovnika
	var jelo *models.Jelo
	for _, j := range append(append(jelovnikPravi.Dorucak, jelovnikPravi.Rucak...), jelovnikPravi.Vecera...) {
		if j.JeloID == jeloID {
			jelo = &j
			break
		}
	}
	if jelo == nil {
		return models.FinansijskaKartica{}, fmt.Errorf("jelo sa ID %s nije pronađeno u jelovniku %s", jeloID, jelovnikID)
	}

	// 4. Parsiraj datum jelovnika
	datumJelovnika, err := time.Parse(time.RFC3339, jelovnikPravi.Datum)
	if err != nil {
		return models.FinansijskaKartica{}, fmt.Errorf("nepravilan format datuma jelovnika: %v", err)
	}

	// 5. Dohvati karticu iz Mongo
	oid, _ := primitive.ObjectIDFromHex(userID)
	var kartica models.FinansijskaKartica
	err = s.Collection.FindOne(context.TODO(), bson.M{"userId": oid}).Decode(&kartica)
	if err != nil {
		return models.FinansijskaKartica{}, err
	}

	// 6. Proveri dnevni limit po tipu obroka (maks 2 po tipu po danu)
	count := 0
	for _, isk := range kartica.IskoriscenaJela {
		if isk.Datum.Equal(datumJelovnika) && isk.TipObroka == jelo.TipObroka {
			count++
		}
	}
	if count >= 2 {
		return models.FinansijskaKartica{}, fmt.Errorf("već ste iskoristili maksimalan broj %s za datum %s", jelo.TipObroka, datumJelovnika.Format("02.01.2006"))
	}

	// 7. Poziv ka jelovnik-servisu da proveri remaining
	remainingURL := fmt.Sprintf("http://host.docker.internal:81/menza/jelovnik/%s/jela/%s/remaining", jelovnikID, jeloID)
	resp, err := http.Get(remainingURL)
	if err != nil {
		return models.FinansijskaKartica{}, fmt.Errorf("greska pri dohvatanju remaining: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.FinansijskaKartica{}, fmt.Errorf("greska: status %d pri dohvatanju remaining", resp.StatusCode)
	}

	var data struct {
		Remaining int `json:"remaining"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.FinansijskaKartica{}, fmt.Errorf("greska pri parsiranju remaining: %v", err)
	}

	if data.Remaining <= 0 {
		return models.FinansijskaKartica{}, fmt.Errorf("nije moguce iskoristiti jelo, limit je dostignut")
	}

	// 8. REST poziv da rezerviše jelo (smanji remaining)
	reserveURL := fmt.Sprintf("http://host.docker.internal:81/menza/jelovnik/%s/jela/%s/reserve", jelovnikID, jeloID)
	resResp, err := http.Post(reserveURL, "application/json", nil)
	if err != nil {
		return models.FinansijskaKartica{}, fmt.Errorf("greska pri rezervaciji jela: %v", err)
	}
	defer resResp.Body.Close()
	if resResp.StatusCode != http.StatusOK {
		return models.FinansijskaKartica{}, fmt.Errorf("nije moguce rezervisati jelo, status %d", resResp.StatusCode)
	}

	if data.Remaining <= 2 {
		notifBody := map[string]interface{}{
			"title":      fmt.Sprintf("Ostatak jela: %s", jelo.Naziv),
			"message":    fmt.Sprintf("Za jelovnik %s, ostalo je još %d porcije jela %s", jelovnikID, data.Remaining, jelo.Naziv),
			"type":       "obrok",
			"jelovnikID": jelovnikID,
			"jeloID":     jeloID,
			"jeloNaziv":  jelo.Naziv,
			"datum":      datumJelovnika.UTC().Format(time.RFC3339), // mora biti "datum" da handler čita
			"remaining":  data.Remaining,
		}

		bodyBytes, _ := json.Marshal(notifBody)
		resp, err := http.Post("http://notification-service:8089/jelo-remaining", "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			fmt.Printf("Greška pri slanju notifikacije: %v\n", err)
		} else {
			defer resp.Body.Close()
			respBody, _ := io.ReadAll(resp.Body)
			fmt.Println("Notification service odgovor:", string(respBody))
		}
	}

	// 10. Smanji count po tipu obroka u kartici
	switch jelo.TipObroka {
	case "dorucak":
		if kartica.DorucakCount <= 0 {
			return models.FinansijskaKartica{}, fmt.Errorf("nemate dovoljno doručaka")
		}
		kartica.DorucakCount--
	case "rucak":
		if kartica.RucakCount <= 0 {
			return models.FinansijskaKartica{}, fmt.Errorf("nemate dovoljno ručkova")
		}
		kartica.RucakCount--
	case "vecera":
		if kartica.VeceraCount <= 0 {
			return models.FinansijskaKartica{}, fmt.Errorf("nemate dovoljno večera")
		}
		kartica.VeceraCount--
	default:
		return models.FinansijskaKartica{}, fmt.Errorf("nepoznat tip obroka")
	}

	// 11. Dodaj u istoriju iskorišćenih jela
	kartica.IskoriscenaJela = append(kartica.IskoriscenaJela, models.IskoriscenoJelo{
		Datum:     datumJelovnika,
		JeloID:    jelo.JeloID,
		Naziv:     jelo.Naziv,
		TipObroka: jelo.TipObroka,
	})

	// 12. Update baze
	_, err = s.Collection.UpdateOne(
		context.TODO(),
		bson.M{"userId": oid},
		bson.M{"$set": kartica},
	)
	if err != nil {
		return models.FinansijskaKartica{}, err
	}

	// 13. Vrati ažuriranu karticu
	return kartica, nil
}

func GetJelovnici() ([]models.Jelovnik, error) {
	resp, err := http.Get("http://host.docker.internal:81/menza/jelovnik/jelovnici-sa-jelima")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("RESPONSE BODY:", string(body)) // <--- ovo će ti pokazati šta stvarno stiže

	var jelovnici []models.Jelovnik
	if err := json.Unmarshal(body, &jelovnici); err != nil {
		return nil, err
	}

	return jelovnici, nil
}
