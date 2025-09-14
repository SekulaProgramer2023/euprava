package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	"review-service/models"
	"review-service/service"
)

// GET /reviews -> svi review-i
func GetAllReviewsHandler(w http.ResponseWriter, r *http.Request) {
	reviews, err := services.GetAllReviews()
	if err != nil {
		http.Error(w, "Failed to fetch reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// GET /reviews/{sobaId} -> review-i za određenu sobu
func GetReviewsBySobaHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	jeloID := params["jeloId"]

	reviews, err := services.GetReviewsByJeloID(jeloID)
	if err != nil {
		http.Error(w, "Failed to fetch reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// POST /reviews -> kreiranje review-a
func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Review

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	created, err := services.CreateReview(input)
	if err != nil {
		http.Error(w, "Failed to create review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func GetAverageRatingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sobaID := params["jeloId"]

	average, err := services.CalculateAverageRating(sobaID)
	if err != nil {
		http.Error(w, "Failed to calculate average rating: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		SobaID  string  `json:"jeloId"`
		Average float64 `json:"average_rating"`
	}{
		SobaID:  sobaID,
		Average: average,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
