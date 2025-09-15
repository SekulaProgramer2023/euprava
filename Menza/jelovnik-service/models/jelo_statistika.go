package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JeloStatistika struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	JelovnikID     primitive.ObjectID `bson:"jelovnikId" json:"jelovnikId"`
	JeloID         primitive.ObjectID `bson:"jeloId" json:"jeloId"`
	BrojPorudzbina int                `bson:"brojPorudzbina" json:"brojPorudzbina"`
	Limit          int                `bson:"limit" json:"limit"`
}
