package main

import (
	"net/http"

	"github.com/Shredder42/learn-http-servers/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, req *http.Request) {
	authorID := req.URL.Query().Get("author_id")
	dbChirps := []database.Chirp{} // not sure about this underline. might be better to use var dbChirps []database.Chirp
	Chirps := []Chirp{}
	var err error

	if len(authorID) > 0 {
		parsedUserID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID", parseErr)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUser(req.Context(), parsedUserID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
			return
		}
	} else {
		dbChirps, err = cfg.db.GetChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
			return
		}
	}

	for _, dbChirp := range dbChirps {
		Chirps = append(Chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, Chirps)

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
