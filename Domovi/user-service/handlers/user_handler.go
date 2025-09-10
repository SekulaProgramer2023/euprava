package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := service.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// RegisterUserHandler obrađuje HTTP request za registraciju korisnika
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
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	token, err := service.LoginUser(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	// parsiraj body
	var request struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// pozovi servis
	user, err := service.GetUserByEmail(request.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// pošalji JSON odgovor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
