package service

import (
	"context"
	"errors"
	"jelovnik-service/db"
	"jelovnik-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Validne kategorije jela
var validKategorije = map[models.KategorijaJela]bool{
	models.Meso:           true,
	models.Vegetarijansko: true,
	models.Kuvano:         true,
	models.Desert:         true,
	models.Predjelo:       true,
	models.Salata:         true,
}

// Validni tipovi obroka
var validTipovi = map[models.TipObroka]bool{
	models.Dorucak: true,
	models.Rucak:   true,
	models.Vecera:  true,
}

func GetJela() ([]models.Jelo, error) {
	collection := db.Client.Database("eupravaM").Collection("jela")
	var jela []models.Jelo

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &jela); err != nil {
		return nil, err
	}

	return jela, nil
}

func CreateJelo(jelo models.Jelo) (*models.Jelo, error) {

	if jelo.Naziv == "" {
		return nil, errors.New("naziv jela ne može biti prazan")
	}

	if _, ok := validKategorije[jelo.Kategorija]; !ok {
		return nil, errors.New("nevažeća kategorija jela")
	}

	if _, ok := validTipovi[jelo.TipObroka]; !ok {
		return nil, errors.New("nevažeći tip obroka")
	}

	if jelo.Kalorije < 0 {
		return nil, errors.New("kalorije ne mogu biti negativne")
	}

	if jelo.Nutritijenti == nil {
		jelo.Nutritijenti = make(map[string]float64)
	}

	collection := db.Client.Database("eupravaM").Collection("jela")
	jelo.JeloID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), jelo)
	if err != nil {
		return nil, err
	}

	return &jelo, nil
}
