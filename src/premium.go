package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	// upgrades a user to premium
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// only accept "user.upgraded" events for now
	if params.Event != "user.upgraded" {
		// if not, respond with 204
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// update the user to premium
	_, err = cfg.db.UpgradeUserToPremium(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User does not exist", err)
		return
	}

	// success - respond with 204
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
