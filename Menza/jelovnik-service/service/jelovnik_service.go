package service

import (
	"context"
	"errors"
	"jelovnik-service/db"
	"jelovnik-service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Kreira jelovnik i razdvaja jela po tipu obroka
func CreateJelovnik(jeloIDs []primitive.ObjectID, opis string, datum time.Time) (*models.Jelovnik, error) {
	if len(jeloIDs) == 0 {
		return nil, errors.New("jelovnik mora imati bar jedno jelo")
	}

	collectionJela := db.Client.Database("eupravaM").Collection("jela")
	var jela []models.Jelo

	// Pronalazi jela iz baze na osnovu ID-ova
	cursor, err := collectionJela.Find(context.TODO(), bson.M{"_id": bson.M{"$in": jeloIDs}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &jela); err != nil {
		return nil, err
	}

	// Proverava da li su pronađena sva jela
	if len(jela) != len(jeloIDs) {
		return nil, errors.New("neka od jela nisu pronađena u bazi")
	}

	dorucakIDs := []primitive.ObjectID{}
	rucakIDs := []primitive.ObjectID{}
	veceraIDs := []primitive.ObjectID{}

	// Razvrstava jela po tipu obroka
	for _, j := range jela {
		switch j.TipObroka {
		case models.Dorucak:
			dorucakIDs = append(dorucakIDs, j.JeloID)
		case models.Rucak:
			rucakIDs = append(rucakIDs, j.JeloID)
		case models.Vecera:
			veceraIDs = append(veceraIDs, j.JeloID)
		default:
			return nil, errors.New("nevažeći tip obroka za jelo: " + j.Naziv)
		}
	}

	// ✅ Validacija – mora imati bar po jedno jelo za svaki tip
	if len(dorucakIDs) == 0 || len(rucakIDs) == 0 || len(veceraIDs) == 0 {
		return nil, errors.New("jelovnik mora imati bar jedan doručak, jedan ručak i jednu večeru")
	}

	jelovnik := models.Jelovnik{
		JelovnikID: primitive.NewObjectID(),
		Datum:      datum,
		Dorucak:    dorucakIDs,
		Rucak:      rucakIDs,
		Vecera:     veceraIDs,
		Opis:       opis,
	}

	// Ubacivanje jelovnika u bazu
	collectionJelovnik := db.Client.Database("eupravaM").Collection("jelovnici")
	_, err = collectionJelovnik.InsertOne(context.TODO(), jelovnik)
	if err != nil {
		return nil, err
	}

	return &jelovnik, nil
}

func GetJelovnike() ([]models.Jelovnik, error) {
	collection := db.Client.Database("eupravaM").Collection("jelovnici")
	var jelovnici []models.Jelovnik

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &jelovnici); err != nil {
		return nil, err
	}

	return jelovnici, nil
}

// Dohvata jelovnike sa imenima jela umesto samo ID-jeva
func GetJelovniciSaJelima() ([]map[string]interface{}, error) {
	jelovnici, err := GetJelovnike()
	if err != nil {
		return nil, err
	}

	// Kolekcija jela
	collectionJela := db.Client.Database("eupravaM").Collection("jela")

	var result []map[string]interface{}

	for _, jelovnik := range jelovnici {
		// Funkcija koja vraća listu jela po ID-jevima
		mapJela := func(ids []primitive.ObjectID) ([]models.Jelo, error) {
			if len(ids) == 0 {
				return []models.Jelo{}, nil
			}
			var jela []models.Jelo
			cursor, err := collectionJela.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				return nil, err
			}
			if err := cursor.All(context.TODO(), &jela); err != nil {
				return nil, err
			}
			return jela, nil
		}

		dorucak, _ := mapJela(jelovnik.Dorucak)
		rucak, _ := mapJela(jelovnik.Rucak)
		vecera, _ := mapJela(jelovnik.Vecera)

		result = append(result, map[string]interface{}{
			"jelovnikId": jelovnik.JelovnikID,
			"datum":      jelovnik.Datum,
			"dorucak":    dorucak,
			"rucak":      rucak,
			"vecera":     vecera,
			"opis":       jelovnik.Opis,
		})
	}

	return result, nil
}
