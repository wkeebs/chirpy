package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32 // thread safe
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	msg := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(msg))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Metrics reset"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	fmt.Println("Starting server...")

	const prefix = "/app"
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux() // router
	config := apiConfig{}

	serveMux.Handle(prefix+"/", config.middlewareMetricsInc(http.StripPrefix(prefix, http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc("/healthz", readinessHandler)
	serveMux.HandleFunc("/metrics", config.metricsHandler)
	serveMux.HandleFunc("/reset", config.resetHandler)

	server := http.Server{Addr: ":" + port, Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
