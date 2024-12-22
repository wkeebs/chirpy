package main

import (
	"fmt"
	"log"
	"net/http"
)



func readinessHandler(w http.ResponseWriter, r *http.Request) {
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
	serveMux.Handle(prefix+"/", http.StripPrefix(prefix, http.FileServer(http.Dir(filepathRoot))))
	serveMux.HandleFunc("/healthz", readinessHandler)

	server := http.Server{Addr: ":" + port, Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
