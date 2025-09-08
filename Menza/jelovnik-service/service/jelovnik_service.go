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
func CreateJelovnik(dorucakIDs, rucakIDs, veceraIDs []primitive.ObjectID, opis string, datum time.Time) (*models.Jelovnik, error) {
	// Validacija
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

	collection := db.Client.Database("eupravaM").Collection("jelovnici")
	_, err := collection.InsertOne(context.TODO(), jelovnik)
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
