package main

import (
	"fmt"
	"jelovnik-service/bootstrap"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"jelovnik-service/db"
	"jelovnik-service/handlers"
)

func main() {
	// Povezivanje na MongoDB
	err := db.ConnectToMongo()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer db.DisconnectMongo()
	bootstrap.ClearJelovnici()
	bootstrap.ClearJela()
	bootstrap.InsertInitialJela()
	bootstrap.InsertInitialJelovnici()
	// Kreiranje routera
	router := mux.NewRouter()

	// Endpoints za jelovnik
	router.HandleFunc("/jela", handlers.GetJela).Methods("GET", "OPTIONS")
	router.HandleFunc("/jela", handlers.CreateJelo).Methods("POST", "OPTIONS")
	router.HandleFunc("/jela/tip", handlers.GetJelaByTipHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/jela/{id}", handlers.GetJeloByIDHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/jelovnik", handlers.CreateJelovnikHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/jelovnik", handlers.GetJelovnikeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/jelovnici-sa-jelima", handlers.GetJelovniciSaJelimaHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/{jelovnikId}/jela/{jeloId}/reserve", handlers.ReserveJeloHandler).Methods("POST")
	router.HandleFunc("/{jelovnikId}/jela/{jeloId}/remaining", handlers.GetRemainingJeloHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/jelovnik/{jelovnikId}", handlers.GetJelovnikByIDHandler).Methods("GET", "OPTIONS")

	// CORS konfiguracija
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // front-end origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Konfiguracija i pokretanje servera
	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8083", // port jelovnik-service
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Jelovnik service started on port 8083")
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting jelovnik service:", err)
		os.Exit(1)
	}
}
