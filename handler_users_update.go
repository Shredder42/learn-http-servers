package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shredder42/learn-http-servers/internal/auth"
	"github.com/Shredder42/learn-http-servers/internal/database"
)

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting bearer token", err)
		return
	}

	UserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid web token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error securing password", err)
	}

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now().UTC(),
		ID:             UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user info", err)
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
