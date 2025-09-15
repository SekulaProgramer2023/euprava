package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification - univerzalan model za obaveštenja
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    string             `bson:"user_id" json:"user_id"`                           // kome je notifikacija namenjena
	Title     string             `bson:"title" json:"title"`                               // kratak naslov notifikacije
	Message   string             `bson:"message" json:"message"`                           // detaljnija poruka
	Type      string             `bson:"type" json:"type"`                                 // tip: "kvar", "review", "dogadjaj", "system", ...
	RelatedId string             `bson:"related_id,omitempty" json:"related_id,omitempty"` // opcioni ID entiteta (npr. sobaId, kvarId, reviewId)
	IsRead    bool               `bson:"is_read" json:"is_read"`                           // da li je korisnik pročitao
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`                     // kada je kreirana
}

// Factory funkcija za kreiranje nove notifikacije
func NewNotification(userId, title, message, notifType, relatedId string) Notification {
	return Notification{
		ID:        primitive.NewObjectID(),
		UserId:    userId,
		Title:     title,
		Message:   message,
		Type:      notifType,
		RelatedId: relatedId,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
}
