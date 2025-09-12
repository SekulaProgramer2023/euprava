package handlers

import (
	"dogadjaj-service/models"
	"dogadjaj-service/service"
	"encoding/json"
	"net/http"
)

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
