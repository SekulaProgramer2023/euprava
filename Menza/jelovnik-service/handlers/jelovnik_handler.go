package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jelovnik-service/service"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateJelovnikRequest struct {
	Dorucak []string `json:"dorucak"`
	Rucak   []string `json:"rucak"`
	Vecera  []string `json:"vecera"`
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

	// Konvertuj ID-eve u ObjectID za svaki tip
	dorucakIDs := []primitive.ObjectID{}
	for _, idStr := range req.Dorucak {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, "Neispravan dorucak ID: "+idStr, http.StatusBadRequest)
			return
		}
		dorucakIDs = append(dorucakIDs, objID)
	}

	rucakIDs := []primitive.ObjectID{}
	for _, idStr := range req.Rucak {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, "Neispravan rucak ID: "+idStr, http.StatusBadRequest)
			return
		}
		rucakIDs = append(rucakIDs, objID)
	}

	veceraIDs := []primitive.ObjectID{}
	for _, idStr := range req.Vecera {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, "Neispravan vecera ID: "+idStr, http.StatusBadRequest)
			return
		}
		veceraIDs = append(veceraIDs, objID)
	}

	jelovnik, err := service.CreateJelovnik(dorucakIDs, rucakIDs, veceraIDs, req.Opis, datum)
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
func GetJelovniciSaJelimaHandler(w http.ResponseWriter, r *http.Request) {
	jelovnici, err := service.GetJelovniciSaJelima()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jelovnici)
}
func ReserveJeloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jelovnikID, _ := primitive.ObjectIDFromHex(vars["jelovnikId"])
	jeloID, _ := primitive.ObjectIDFromHex(vars["jeloId"])

	err := service.ReserveJelo(jelovnikID, jeloID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reserved"))
}
func GetRemainingJeloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jelovnikIDHex := vars["jelovnikId"]
	jeloIDHex := vars["jeloId"]

	jelovnikID, err := primitive.ObjectIDFromHex(jelovnikIDHex)
	if err != nil {
		http.Error(w, "nepravilan jelovnikID", http.StatusBadRequest)
		return
	}

	jeloID, err := primitive.ObjectIDFromHex(jeloIDHex)
	if err != nil {
		http.Error(w, "nepravilan jeloID", http.StatusBadRequest)
		return
	}

	remaining, err := service.GetRemaining(jelovnikID, jeloID)
	if err != nil {
		http.Error(w, "greska pri dohvatanju remaining", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"remaining": remaining,
	})
}
func GetJelovnikByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jelovnikIDHex := vars["jelovnikId"]

	jelovnikID, err := primitive.ObjectIDFromHex(jelovnikIDHex)
	if err != nil {
		http.Error(w, "Neispravan jelovnikID", http.StatusBadRequest)
		return
	}

	jelovnik, err := service.GetJelovnikByID(jelovnikID)
	if err != nil {
		http.Error(w, "Jelovnik nije pronađen", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jelovnik)
}
