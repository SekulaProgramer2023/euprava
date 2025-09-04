package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Soba struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomNumber string             `bson:"roomNumber" json:"roomNumber"`
	Capacity   int                `bson:"capacity" json:"capacity"`
	Users      []string           `bson:"users" json:"users"`
	OnBudget   bool               `bson:"onBudget" json:"onBudget"`
	IsFree     bool               `bson:"IsFree" json:"IsFree"`
}

func NewSoba(roomNumber string, capacity int, onBudget, isFree bool) Soba {
	return Soba{
		ID:         primitive.NewObjectID(),
		RoomNumber: roomNumber,
		Capacity:   capacity,
		Users:      []string{},
		OnBudget:   onBudget,
		IsFree:     isFree,
	}
}
