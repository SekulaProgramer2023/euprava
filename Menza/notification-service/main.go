package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"notification-service/bootstrap"
	"notification-service/db"
	"notification-service/handlers"
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
	bootstrap.ClearNotifications()

	router := mux.NewRouter()
	router.HandleFunc("/jelo-remaining", handlers.CreateJeloRemainingNotificationHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/notification", handlers.GetAllNotificationsHandler).Methods("GET", "OPTIONS")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8089",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Notification service started on port 8089")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting notification service:", err)
		os.Exit(1)
	}
}
