package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FinansijskaKartica struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	Ime             string             `bson:"ime" json:"ime"`
	Prezime         string             `bson:"prezime" json:"prezime"`
	Index           string             `bson:"index" json:"index"`
	Novac           float64            `bson:"novac" json:"novac"`
	DorucakCount    int                `bson:"dorucakCount" json:"dorucakCount"`
	RucakCount      int                `bson:"rucakCount" json:"rucakCount"`
	VeceraCount     int                `bson:"veceraCount" json:"veceraCount"`
	IskoriscenaJela []IskoriscenoJelo  `bson:"iskoriscenaJela" json:"iskoriscenaJela"`
}
type IskoriscenoJelo struct {
	Datum     time.Time `bson:"datum" json:"datum"`
	JeloID    string    `bson:"jeloId" json:"jeloId"`
	Naziv     string    `bson:"naziv" json:"naziv"`
	TipObroka string    `bson:"tipObroka" json:"tipObroka"`
}

func NewFinansijskaKartica(userID primitive.ObjectID, ime, prezime, index string) FinansijskaKartica {
	return FinansijskaKartica{
		UserID:       userID,
		Ime:          ime,
		Prezime:      prezime,
		Index:        index,
		Novac:        0.0,
		DorucakCount: 0,
		RucakCount:   0,
		VeceraCount:  0,
	}
}
