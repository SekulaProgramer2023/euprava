package handlers

import (
	"encoding/json"
	"net/http"
	"sobe-service/models"
	"sobe-service/service"
	"strings"

	"github.com/gorilla/mux"
)

// ✅ Kreiranje novog kvara (POST /kvarovi)
func CreateKvarHandler(w http.ResponseWriter, r *http.Request) {
	var kvar models.Kvar
	if err := json.NewDecoder(r.Body).Decode(&kvar); err != nil {
		http.Error(w, "Nevalidan input", http.StatusBadRequest)
		return
	}

	newKvar := models.NewKvar(kvar.UserId, kvar.SobaId, kvar.Description, kvar.Status)

	createdKvar, err := services.CreateKvar(newKvar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdKvar)
}

// ✅ Lista svih kvarova (GET /kvarovi)
func GetAllKvaroviHandler(w http.ResponseWriter, r *http.Request) {
	kvarovi, err := services.GetAllKvarovi()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kvarovi)
}

// ✅ Lista kvarova po sobi (GET /kvarovi/soba/{id})
func GetKvaroviBySobaHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sobaID := params["id"]

	kvarovi, err := services.GetKvaroviBySobaID(sobaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kvarovi)
}

func ResolveKvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	kvarID := strings.TrimSpace(vars["id"])

	if kvarID == "" {
		http.Error(w, `{"error": "Nedostaje ID kvara"}`, http.StatusBadRequest)
		return
	}

	err := services.ResolveKvar(kvarID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kvar uspešno označen kao rešen",
	})
}
