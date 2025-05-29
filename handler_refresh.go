package main

import (
	"net/http"
	"time"

	"github.com/Shredder42/learn-http-servers/internal/auth"
)

type Token struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting bearer token", err)
		return
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired", err)
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token was revoked", err)
		return
	}

	tokenString, err := auth.MakeJWT(refreshToken.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get token string", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Token{
		Token: tokenString,
	})

}
