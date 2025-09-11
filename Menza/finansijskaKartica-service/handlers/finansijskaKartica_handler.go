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

// POST /kartice/{userId}/deposit
func (h *KarticaHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDHex := vars["userId"]

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	var req struct {
		Novac float64 `json:"novac"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.Deposit(userID, req.Novac)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
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

// Kupovina doručka
func (h *KarticaHandler) BuyRuckoviHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDHex := vars["userId"]

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	var req struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.BuyRuckovi(userID, req.Count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *KarticaHandler) BuyVecereHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDHex := vars["userId"]

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	var req struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.BuyVecere(userID, req.Count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
func (h *KarticaHandler) BuyDorucakHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDHex := vars["userId"]

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	var req struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.BuyDorucak(userID, req.Count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
func (h *KarticaHandler) IskoristiObrokHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	jelovnikID := vars["jelovnikId"] // novo
	jeloID := vars["jeloId"]

	kartica, err := h.Service.IskoristiObrok(userID, jelovnikID, jeloID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kartica)
}
