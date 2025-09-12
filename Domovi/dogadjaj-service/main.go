package main

import (
	"dogadjaj-service/bootstrap"
	"dogadjaj-service/db"
	"dogadjaj-service/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"time"
)

func main() {
	err := db.ConnectToMongo()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer db.DisconnectMongo()
	bootstrap.ClearDogadjaj()
	bootstrap.InsertInitialDogadjaji()

	router := mux.NewRouter()
	router.HandleFunc("/dogadjaji", handlers.GetAllDogadjajiHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/dogadjaj", handlers.CreateDogadjajHandler).Methods("POST", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8084",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Dogadjaj service started on port 8084")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting user service:", err)
		os.Exit(1)
	}
}
