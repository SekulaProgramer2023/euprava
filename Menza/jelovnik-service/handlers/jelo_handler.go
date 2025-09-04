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
func GetJelaByTipHandler(w http.ResponseWriter, r *http.Request) {
	tip := r.URL.Query().Get("tip")
	if tip == "" {
		http.Error(w, "parametar 'tip' je obavezan", http.StatusBadRequest)
		return
	}

	// Konverzija u models.TipObroka
	tipObroka := models.TipObroka(tip)

	jela, err := service.GetJelaByTip(tipObroka)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jela)
}
