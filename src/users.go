package main

import (
	"encoding/json"
	"net/http"

	"github.com/wkeebs/chirpy/internal/auth"
	"github.com/wkeebs/chirpy/internal/database"
)

func (cfg *apiConfig) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.GetAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get Chirps", err)
		return
	}

	// map for correct json representation
	var respUsers []User
	for _, u := range users {
		respUsers = append(respUsers, User{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Email:     u.Email,
		})
	}

	// write response
	respondWithJSON(w, http.StatusOK, respUsers)
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating User", nil)
		return
	}

	// create user
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create User", nil)
		return
	}

	// map from databaser user struct
	respUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	// construct response
	respondWithJSON(w, http.StatusCreated, respUser)
}
