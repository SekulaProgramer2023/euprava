package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Dogadjaj struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Naziv              string             `bson:"naziv" json:"naziv"`
	Opis               string             `bson:"opis" json:"opis"`
	DatumOdrzavanja    time.Time          `bson:"datum_odrzavanja" json:"datum_odrzavanja"`
	DatumSlanjaZahteva time.Time          `bson:"datum_slanja_zahteva" json:"datum_slanja_zahteva"`
	DatumOdgovora      *time.Time         `bson:"datum_odgovora,omitempty" json:"datum_odgovora,omitempty"`
	Tema               string             `bson:"tema" json:"tema"`
	Users              []string           `bson:"users" json:"users"`
	Status             string             `bson:"status" json:"status"` // "prihvaćen" | "odbijen" | "na čekanju"
}

// Konstruktor
func NewDogadjaj(naziv, opis string, datumOdrzavanja time.Time, tema string) Dogadjaj {
	return Dogadjaj{
		ID:                 primitive.NewObjectID(),
		Naziv:              naziv,
		Opis:               opis,
		DatumOdrzavanja:    datumOdrzavanja,
		DatumSlanjaZahteva: time.Now(),
		Status:             "na čekanju", // default
		Tema:               tema,
		Users:              []string{},
	}
}
