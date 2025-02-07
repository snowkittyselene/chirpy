package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const port = "8080"
const rootFilePath = "."

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerCountRequests(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!", cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func main() {
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath)))
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidation)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCountRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s", rootFilePath, port)
	log.Fatal(srv.ListenAndServe())
}
