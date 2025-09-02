package handlers

import (
	"encoding/json"
	"jelovnik-service/service"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateJelovnikRequest struct {
	JeloIDs []string `json:"jela"`
	Opis    string   `json:"opis,omitempty"`
	Datum   string   `json:"datum"` // u formatu "YYYY-MM-DD"
}

func CreateJelovnikHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateJelovnikRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Neispravan zahtev", http.StatusBadRequest)
		return
	}

	// Parse datum
	datum, err := time.Parse("2006-01-02", req.Datum)
	if err != nil {
		http.Error(w, "Neispravan format datuma", http.StatusBadRequest)
		return
	}

	// Konvertuj string ID-eve u ObjectID
	jeloIDs := []primitive.ObjectID{}
	for _, idStr := range req.JeloIDs {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, "Neispravan JeloID: "+idStr, http.StatusBadRequest)
			return
		}
		jeloIDs = append(jeloIDs, objID)
	}

	jelovnik, err := service.CreateJelovnik(jeloIDs, req.Opis, datum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jelovnik)
}

func GetJelovnikeHandler(w http.ResponseWriter, r *http.Request) {
	jelovnici, err := service.GetJelovnike()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jelovnici)
}
