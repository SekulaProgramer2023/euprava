package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"sobe-service/bootstrap"
	"time"

	"sobe-service/db"
	"sobe-service/handlers"
)

func main() {
	err := db.ConnectToMongo()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer db.DisconnectMongo()
	bootstrap.ClearUsers()
	bootstrap.InsertInitialSobe()

	router := mux.NewRouter()
	router.HandleFunc("/sobe", handlers.GetSobeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/sobe/{id}", handlers.GetSobaByIDHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/sobeSaKapacitetom", handlers.GetSobeWithCapacityHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/kreirajSobu", handlers.CreateSobaHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/useliStudenta", handlers.UseliUseraHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/kvarovi", handlers.GetAllKvaroviHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/prijavi-kvar", handlers.CreateKvarHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/kvarovi/soba/{id}", handlers.GetKvaroviBySobaHandler).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8082",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("User service started on port 8082")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting user service:", err)
		os.Exit(1)
	}
}
