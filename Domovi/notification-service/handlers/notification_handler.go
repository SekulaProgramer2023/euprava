package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"notification-service/service"

	"github.com/gorilla/mux"
)

// POST /notifications
func CreateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notif map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := services.CreateNotification(context.Background(), notif)
	if err != nil {
		http.Error(w, "Error creating notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification created"})
}

// GET /notifications
func GetAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	results, err := services.GetAllNotifications(context.Background())
	if err != nil {
		http.Error(w, "Error fetching notifications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GET /notifications/user/{id}
func GetNotificationsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	results, err := services.GetNotificationsByUser(context.Background(), userId)
	if err != nil {
		http.Error(w, "Error fetching notifications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
