package main

import (
	"log"
	"net/http"
)

const port = "8080"
const rootFilePath = "."

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(rootFilePath)))
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s", rootFilePath, port)
	log.Fatal(srv.ListenAndServe())
}
