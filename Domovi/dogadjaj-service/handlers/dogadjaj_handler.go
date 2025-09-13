package handlers

import (
	"dogadjaj-service/models"
	"dogadjaj-service/service"
	"encoding/json"
	"net/http"
	"strings"
)

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// POST /dogadjaji
func CreateDogadjajHandler(w http.ResponseWriter, r *http.Request) {
	var dogadjaj models.Dogadjaj
	if err := json.NewDecoder(r.Body).Decode(&dogadjaj); err != nil {
		http.Error(w, "neispravan JSON", http.StatusBadRequest)
		return
	}

	created, err := services.CreateDogadjaj(dogadjaj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GET /dogadjaji
func GetAllDogadjajiHandler(w http.ResponseWriter, r *http.Request) {
	dogadjaji, err := services.GetAllDogadjaji()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dogadjaji)
}

func UpdateDogadjajStatusHandler(w http.ResponseWriter, r *http.Request) {
	// očekujemo URL tipa: /dogadjaji/<id>/status
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "nedostaje ID događaja", http.StatusBadRequest)
		return
	}
	id := parts[2]

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "neispravan JSON", http.StatusBadRequest)
		return
	}

	if err := services.UpdateDogadjajStatus(id, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "status uspešno ažuriran"})
}
