package main

import (
	"log"
	"net/http"
)

const PORT = "8080"

func main() {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}
	log.Printf("Server running on port: %s", PORT)
	log.Fatal(srv.ListenAndServe())
}
