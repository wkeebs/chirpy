package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	// forbidden outside of dev environment
	if cfg.platform != "dev" {
		respondWithJSON(w, http.StatusForbidden, "Access denied outside of dev enviroment")
		return
	}

	// delete all users
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete users", err)
		return
	}

	// write response
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Metrics reset"))

	log.Printf("Reset complete")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := cfg.fileserverHits.Load()
	htmlData, err := os.ReadFile("metrics.html")
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf(string(htmlData), hits)
	w.Write([]byte(msg))
}
