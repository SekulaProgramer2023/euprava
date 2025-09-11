package main

import (
	"finansijskaKartica-service/bootstrap"
	"finansijskaKartica-service/handlers"
	"finansijskaKartica-service/service"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"time"

	"finansijskaKartica-service/db"
)

func main() {
	err := db.ConnectToMongo()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer db.DisconnectMongo()
	bootstrap.ClearKartice()
	dbInstance := db.Client.Database("eupravaM")
	karticaService := service.NewFinansijskaKarticaService(dbInstance)
	karticaHandler := handlers.KarticaHandler{Service: karticaService}

	router := mux.NewRouter()
	router.HandleFunc("/kartice", karticaHandler.CreateKarticaHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/kartice", karticaHandler.GetKarticeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/kartice/user/{userId}", karticaHandler.GetKarticaByUserHandler).Methods("GET")
	router.HandleFunc("/kartice/{userId}/deposit", karticaHandler.DepositHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/kartice/{userId}/buy/rucak", karticaHandler.BuyRuckoviHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/kartice/{userId}/buy/vecera", karticaHandler.BuyVecereHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/kartice/{userId}/buy/dorucak", karticaHandler.BuyDorucakHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/kartice/iskoristi/{userId}/{jelovnikId}/{jeloId}",
		karticaHandler.IskoristiObrokHandler).Methods("POST", "OPTIONS")

	// dodaj druge rute kasnije ako bude potrebno, npr GET po userId

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Server
	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8085",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("FinansijskaKartica service started on port 8082")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting finansijskakartica service:", err)
		os.Exit(1)
	}
}
