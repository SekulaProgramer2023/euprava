package handlers

import (
	"encoding/json"
	"finansijskaKartica-service/models"
	"finansijskaKartica-service/service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KarticaHandler struct {
	Service *service.FinansijskaKarticaService
}

// Kreira karticu preko HTTP POST
func (h *KarticaHandler) CreateKarticaHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"userId"`
		Ime     string `json:"ime"`
		Prezime string `json:"prezime"`
		Index   string `json:"index"`
	}

	fmt.Println("Uspeo")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	kartica := models.NewFinansijskaKartica(userID, req.Ime, req.Prezime, req.Index)

	created, err := h.Service.CreateKartica(kartica)
	if err != nil {
		http.Error(w, "failed to insert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
func (h *KarticaHandler) GetKarticaByUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDHex := vars["userId"]

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	kartica, err := h.Service.GetKarticaByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kartica)
}

// GET /kartice
func (h *KarticaHandler) GetKarticeHandler(w http.ResponseWriter, r *http.Request) {
	kartice, err := h.Service.GetKartice()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kartice)
}
