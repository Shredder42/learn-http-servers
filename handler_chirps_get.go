package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbchirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
	}

	newChirps := []Chirp{}
	for _, dbchirp := range dbchirps {
		newChirps = append(newChirps, Chirp{
			ID:        dbchirp.ID,
			CreatedAt: dbchirp.CreatedAt,
			UpdatedAt: dbchirp.UpdatedAt,
			Body:      dbchirp.Body,
			UserID:    dbchirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, newChirps)

}
