package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/models"
	"user-service/service"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := service.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Dekodiraj JSON iz body-a
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pozovi servis
	createdUser, err := service.RegisterUser(user)
	if err != nil {
		// Ako već postoji korisnik, vrati 409 Conflict
		if err.Error() == fmt.Sprintf("user with email %s already exists", user.Email) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Uspešna registracija
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authUser, err := service.LoginUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authUser)
}
