package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Jelovnik struct {
	JelovnikID primitive.ObjectID   `bson:"_id,omitempty" json:"jelovnikId"`
	Datum      time.Time            `bson:"datum" json:"datum"`
	Dorucak    []primitive.ObjectID `bson:"dorucak,omitempty" json:"dorucak,omitempty"`
	Rucak      []primitive.ObjectID `bson:"rucak,omitempty" json:"rucak,omitempty"`
	Vecera     []primitive.ObjectID `bson:"vecera,omitempty" json:"vecera,omitempty"`
	Opis       string               `bson:"opis,omitempty" json:"opis,omitempty"`
}
