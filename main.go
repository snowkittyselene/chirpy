package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/snowkittyselene/chirpy/internal/database"
)

const port = "8080"
const rootFilePath = "."

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	apiKey         string
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
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("TOKEN_SECRET")
	apiKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath)))
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secret,
		apiKey:         apiKey,
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.HandleFunc("POST /api/users", apiCfg.handlerAddUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerAddChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCountRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s", rootFilePath, port)
	log.Fatal(srv.ListenAndServe())
}
