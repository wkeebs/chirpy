package main

import (
	"net/http"
	"time"

	"github.com/wkeebs/chirpy/internal/auth"
)

// refreshHandler - [POST /api/refresh] : generates a new access token
func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// this endpoint takes no body, but expects an refresh token in the auth header
	refreshTok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token present", err)
		return
	}

	// look up token
	storedToken, err := cfg.db.GetRefreshToken(r.Context(), refreshTok)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// check if the token has been revoked
	if storedToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked", err)
		return
	}

	// check if the token has expired
	if storedToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", err)
		return
	}

	// create new access token for the user
	user, err := cfg.db.GetUserByID(r.Context(), storedToken.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User does not exist", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(time.Hour))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new access token", err)
		return
	}

	// success!
	respondWithJSON(w, http.StatusOK, response{
		Token: newAccessToken,
	})
}

// revokeHandler - [POST /api/revoke] : revokes a user's refresh token
func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	// this endpoint takes no body, but expects an refresh token in the auth header
	refreshTok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token present", err)
		return
	}

	// look up token
	_, err = cfg.db.GetRefreshToken(r.Context(), refreshTok)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// revoke token
	_, err = cfg.db.RevokeToken(r.Context(), refreshTok)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}

	// success - respond with 204
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
