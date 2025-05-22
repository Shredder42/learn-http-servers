package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in dev environment"))
		return
	}
	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
	w.Write([]byte("Cleared users from database"))

}
