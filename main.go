package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux() // router
	serveMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := http.Server{Addr: ":8080", Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
