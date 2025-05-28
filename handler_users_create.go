package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shredder42/learn-http-servers/internal/auth"
	"github.com/Shredder42/learn-http-servers/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Token     string    `json:"token"`
}
type parameters struct {
	ExpiresInSeconds int    `json:"expires_in_seconds"`
	Password         string `json:"password"`
	Email            string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error securing password", err)
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		HashedPassword: hashed_password,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

	fmt.Printf("user with email %s created at %v\n", user.Email, user.CreatedAt)

}
