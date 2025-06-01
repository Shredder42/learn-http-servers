package main

import (
	"net/http"
	"sort"

	"github.com/Shredder42/learn-http-servers/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, req *http.Request) {
	authorID := req.URL.Query().Get("author_id")
	sortMethod := req.URL.Query().Get("sort")
	dbChirps := []database.Chirp{}
	chirps := []Chirp{}
	var err error // this is an interface so err := error{} doesn't work

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
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	// sort direction is asc by default
	if sortMethod != "asc" && sortMethod != "desc" && sortMethod != "" {
		respondWithError(w, http.StatusBadRequest, "invalid sort method", err)
		return
	}

	if sortMethod == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)

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
