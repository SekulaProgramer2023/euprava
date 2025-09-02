package handlers

import (
	"encoding/json"
	"jelovnik-service/models"
	"jelovnik-service/service"
	"net/http"
)

// GET /jela
func GetJela(w http.ResponseWriter, r *http.Request) {
	jela, err := service.GetJela()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jela)
}

// POST /jela
func CreateJelo(w http.ResponseWriter, r *http.Request) {
	var jelo models.Jelo

	if err := json.NewDecoder(r.Body).Decode(&jelo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdJelo, err := service.CreateJelo(jelo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdJelo)
}
