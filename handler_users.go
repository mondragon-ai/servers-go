package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	
	// Struct for decoding incoming request body
	type UserEmail struct {
		Email string `json:"email"`
	}

	// Decode the incoming request body
	decoder := json.NewDecoder(r.Body)
	var email UserEmail
	if err := decoder.Decode(&email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Create a new user in the database
	newUser, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Map the database result to the API's User struct
	responseUser := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	// Respond with the newly created user
	respondWithJSON(w, http.StatusCreated, responseUser)
}
