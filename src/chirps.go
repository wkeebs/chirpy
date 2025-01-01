package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/wkeebs/chirpy/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// validate chirp
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// check user exists
	user, err := cfg.db.GetUser(r.Context(), params.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to find user", err)
		return
	}
	if user.ID != params.UserID {
		respondWithJSON(w, http.StatusNotFound, "User not found")
		return
	}

	cleanedBody := replaceProfanity(params.Body)

	// add to database
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	})

	// create response
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func replaceProfanity(s string) string {
	const blur string = "****"
	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	splitString := strings.Split(s, " ")
	for i, word := range splitString {
		for _, pWord := range profaneWords {
			if strings.ToLower(word) == pWord {
				splitString[i] = blur
			}
		}
	}

	return strings.Join(splitString, " ")
}
