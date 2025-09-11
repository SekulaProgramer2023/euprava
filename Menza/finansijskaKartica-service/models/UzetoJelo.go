package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UzetoJelo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	KarticaID     primitive.ObjectID `bson:"karticaId" json:"karticaId"`
	JeloID        primitive.ObjectID `bson:"jeloId" json:"jeloId"`
	JelovnikID    primitive.ObjectID `bson:"jelovnikId" json:"jelovnikId"`
	TipObroka     string             `bson:"tipObroka" json:"tipObroka"` // dorucak / rucak / vecera
	VremeUzimanja time.Time          `bson:"vremeUzimanja" json:"vremeUzimanja"`
}
