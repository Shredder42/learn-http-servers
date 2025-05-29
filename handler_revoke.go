package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Shredder42/learn-http-servers/internal/auth"
	"github.com/Shredder42/learn-http-servers/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting bearer token", err)
		return
	}

	_, err = cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	}

	nullableTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: nullableTime,
		UpdatedAt: time.Now().UTC(),
		Token:     token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token", err)
	}

	w.WriteHeader(http.StatusNoContent)

	fmt.Println("Refresh Token Revoked")
}
