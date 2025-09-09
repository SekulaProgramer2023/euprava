package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sobe-service/models"
	"sobe-service/service"
)

func GetSobeHandler(w http.ResponseWriter, r *http.Request) {
	sobe, err := services.GetAllSobe()
	if err != nil {
		http.Error(w, "Failed to fetch sobe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sobe)
}

// GET /sobe/capacity -> vraća sobe sa Capacity > 0
func GetSobeWithCapacityHandler(w http.ResponseWriter, r *http.Request) {
	sobe, err := services.GetSobeWithCapacity()
	if err != nil {
		http.Error(w, "Failed to fetch sobe with capacity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sobe)
}

func CreateSobaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RoomNumber string `json:"roomNumber"`
		Capacity   int    `json:"capacity"`
		OnBudget   bool   `json:"onBudget"`
		IsFree     bool   `json:"isFree"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.RoomNumber == "" || input.Capacity <= 0 {
		http.Error(w, "RoomNumber is required and Capacity must be > 0", http.StatusBadRequest)
		return
	}

	soba := models.NewSoba(input.RoomNumber, input.Capacity, input.OnBudget, input.IsFree)
	created, err := services.CreateSoba(soba)
	if err != nil {
		http.Error(w, "Failed to create soba: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func UseliUseraHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RoomID string `json:"roomId"`
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedSoba, err := services.UseliUsera(input.RoomID, input.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSoba)
}

func GetSobaByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	roomID := params["id"]

	soba, err := services.GetSobaByID(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(soba)
}
