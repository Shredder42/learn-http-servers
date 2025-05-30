package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Shredder42/learn-http-servers/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, req *http.Request) {

	type Data struct {
		UserID string `json:"user_id"` // could have set this to UUID as well
	}
	type Webhook struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	webhook := Webhook{}
	err := decoder.Decode(&webhook)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsedUserID, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), database.UpgradeUserParams{
		UpdatedAt: time.Now().UTC(),
		ID:        parsedUserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
	}

	w.WriteHeader(http.StatusNoContent)

	fmt.Println("user upgraded")

}
