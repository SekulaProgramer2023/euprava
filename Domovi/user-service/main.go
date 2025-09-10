package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"time"
	"user-service/bootstrap"

	"user-service/db"
	"user-service/handlers"
)

func main() {
	err := db.ConnectToMongo()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer db.DisconnectMongo()
	bootstrap.ClearUsers()
	bootstrap.InsertInitialUsers()

	router := mux.NewRouter()
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{id}", handlers.GetUserByIDHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/by-email", handlers.GetUserByEmailHandler).Methods("POST", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("User service started on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting user service:", err)
		os.Exit(1)
	}
}
