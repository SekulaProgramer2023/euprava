package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SobaId  string             `bson:"soba_id" json:"soba_id"`
	UserId  string             `bson:"user_id" json:"user_id"`
	Rating  int                `bson:"rating" json:"rating"`   // 1–5
	Comment string             `bson:"comment" json:"comment"` // opcioni tekstualni komentar
}

func NewReview(sobaId string, userId string, rating int, comment string) Review {
	return Review{
		ID:      primitive.NewObjectID(),
		SobaId:  sobaId,
		UserId:  userId,
		Rating:  rating,
		Comment: comment,
	}
}
