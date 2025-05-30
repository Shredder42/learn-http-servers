package main

import (
	"fmt"
	"net/http"

	"github.com/Shredder42/learn-http-servers/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirps(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid web token", err)
		return
	}

	chirpID := req.PathValue("chirpID")

	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "user not authorized to delete chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Println("chirp deleted")

}
