package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"review-service/bootstrap"
	"review-service/db"
	"review-service/handlers"
	"time"
)

func main() {
	err := db.ConnectToMongo()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}

	defer db.DisconnectMongo()
	bootstrap.ClearReviews()

	router := mux.NewRouter()
	router.HandleFunc("/reviews", handlers.GetAllReviewsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/create", handlers.CreateReviewHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/{jeloId}", handlers.GetReviewsBySobaHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/average/{jeloId}", handlers.GetAverageRatingHandler).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8087",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Review service started on port 8087")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting review service:", err)
		os.Exit(1)
	}
}
