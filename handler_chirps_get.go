package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
	}

	newChirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		newChirps = append(newChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, newChirps)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})

}
