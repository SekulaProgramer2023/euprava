package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"notification-service/service"
)

// POST /notifications
// Request payload za kreiranje notifikacije za remaining jelo
type JeloRemainingNotificationRequest struct {
	JelovnikID string `json:"jelovnikId"`
	JeloID     string `json:"jeloId"`
	JeloNaziv  string `json:"jeloNaziv"`
	Remaining  int    `json:"remaining"`
}

// / POST /notifications/jelo-remaining
func CreateJeloRemainingNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JelovnikID    string `json:"jelovnikID"`
		JelovnikNaziv string `json:"jelovnikNaziv"`
		JeloID        string `json:"jeloID"`
		JeloNaziv     string `json:"jeloNaziv"`
		Remaining     int    `json:"remaining"`
		Datum         string `json:"datum"` // RFC3339 format
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parsiranje datuma u RFC3339 formatu
	datum, err := time.Parse(time.RFC3339, req.Datum)
	if err != nil {
		http.Error(w, "Invalid datum format, expected RFC3339: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Kreiranje notifikacije
	err = services.CreateJeloRemainingNotification(
		context.Background(),
		req.JelovnikID,
		req.JelovnikNaziv,
		req.JeloID,
		req.JeloNaziv,
		req.Remaining,
		datum,
	)
	if err != nil {
		http.Error(w, "Error creating notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification created"})
}

func GetAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	// Poziv servisa
	notifications, err := services.GetAllNotifications(context.Background())
	if err != nil {
		http.Error(w, "Greška pri dohvatanju notifikacija: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Vraćanje JSON odgovora
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
