package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Kvar struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId      string             `bson:"user_id" json:"user_id"`
	SobaId      string             `bson:"soba_id" json:"soba_id"`
	Description string             `bson:"description" json:"description"`
	Status      bool               `bson:"status" json:"status"`
}

func NewKvar(userId string, sobaId string, description string, status bool) Kvar {
	return Kvar{
		ID:          primitive.NewObjectID(),
		UserId:      userId,
		SobaId:      sobaId,
		Description: description,
		Status:      status,
	}
}
