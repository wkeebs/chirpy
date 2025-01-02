package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/wkeebs/chirpy/internal/auth"
	"github.com/wkeebs/chirpy/internal/database"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	// unpack chirp id
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid Chirp ID", err)
		return
	}

	// get chirp
	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get Chirp", err)
		return
	}

	// write response
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get Chirps", err)
		return
	}

	// map for correct json representation
	var respChirps []Chirp
	for _, c := range chirps {
		respChirps = append(respChirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	// write response
	respondWithJSON(w, http.StatusOK, respChirps)
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// get JWT from headers
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	// unpack user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to find user", err)
		return
	}
	if user.ID != userID {
		respondWithJSON(w, http.StatusNotFound, "User not found")
		return
	}

	cleanedBody := replaceProfanity(params.Body)

	// add to database
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
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

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	// expects:
	// 1. an access token in the header
	// 2. the ID of the chirp to delete in the path

	// check access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find access token", err)
		return
	}

	// unpack user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// unpack chirp id
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid Chirp ID", err)
		return
	}

	// check that the chirp exists and is authored by the user
	storedChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp does not exist", err)
		return
	}
	if storedChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User is not the author of the chirp", err)
		return
	}

	// delete the chirp
	cfg.db.DeleteChirp(r.Context(), chirpId)

	// verify it has been deleted
	_, err = cfg.db.GetChirp(r.Context(), chirpId)
	if err == nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp was not deleted correctly", err)
		return
	}

	// success - respond with 204
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
