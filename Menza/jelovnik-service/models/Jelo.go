package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Enumeracija za kategoriju jela
type KategorijaJela string

const (
	Meso           KategorijaJela = "meso"
	Vegetarijansko KategorijaJela = "vegetarijansko"
	Kuvano         KategorijaJela = "kuvano"
	Desert         KategorijaJela = "desert"
	Predjelo       KategorijaJela = "predjelo"
	Salata         KategorijaJela = "salata"
)

// Enumeracija za tip obroka
type TipObroka string

const (
	Dorucak TipObroka = "dorucak"
	Rucak   TipObroka = "rucak"
	Vecera  TipObroka = "vecera"
)

// Struktura Jelo
type Jelo struct {
	JeloID       primitive.ObjectID `bson:"_id,omitempty" json:"jeloId"`                          // MongoDB ObjectID
	Naziv        string             `bson:"naziv" json:"naziv"`                                   // naziv jela
	Kategorija   KategorijaJela     `bson:"kategorija" json:"kategorija"`                         // tip jela (meso, vegetarijansko, kuvano, desert)
	TipObroka    TipObroka          `bson:"tipObroka" json:"tipObroka"`                           // dorucak, rucak, vecera, uzina
	Kalorije     int                `bson:"kalorije" json:"kalorije"`                             // kalorije
	Nutritijenti map[string]float64 `bson:"nutritijenti,omitempty" json:"nutritijenti,omitempty"` // npr. {"proteini": 10, "masti": 5}
}
