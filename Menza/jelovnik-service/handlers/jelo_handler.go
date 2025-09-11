package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jelovnik-service/db"
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
func GetJeloByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jeloID := vars["id"]

	oid, err := primitive.ObjectIDFromHex(jeloID)
	if err != nil {
		http.Error(w, "Invalid jelo ID", http.StatusBadRequest)
		return
	}

	collection := db.Client.Database("eupravaM").Collection("jela")
	var jelo models.Jelo
	err = collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&jelo)
	if err != nil {
		http.Error(w, "Jelo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jelo)
}
