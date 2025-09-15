package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification - univerzalan model za obaveštenja
type Notification struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Message        string             `bson:"message" json:"message"`
	Type           string             `bson:"type" json:"type"`
	JelovnikID     string             `bson:"jelovnik_id,omitempty" json:"jelovnik_id,omitempty"`
	JelovnikNaziv  string             `bson:"jelovnik_naziv,omitempty" json:"jelovnik_naziv,omitempty"`
	JeloID         string             `bson:"jelo_id,omitempty" json:"jelo_id,omitempty"`
	JeloNaziv      string             `bson:"jelo_naziv,omitempty" json:"jelo_naziv,omitempty"`
	DatumJelovnika string             `bson:"datum_jelovnika,omitempty" json:"datum_jelovnika,omitempty"` // novo
	IsRead         bool               `bson:"is_read" json:"is_read"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}

// Factory funkcija za kreiranje nove notifikacije
func NewNotification(title, message, notifType string, jelovnikID, jeloID string) Notification {
	return Notification{
		ID:         primitive.NewObjectID(),
		Title:      title,
		Message:    message,
		Type:       notifType,
		JelovnikID: jelovnikID,
		JeloID:     jeloID,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
}
