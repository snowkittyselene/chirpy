package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/snowkittyselene/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}{}
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}
	if len(req.Body) > 140 {
		respondError(w, http.StatusBadRequest, "Chirp is too long", nil)
	}
	newChirp, err := cfg.db.AddChirp(r.Context(), database.AddChirpParams{
		Body:   removeBadWords(req.Body),
		UserID: req.UserID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error adding Chirp to database", err)
		return
	}
	chirp := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving Chirps", err)
		return
	}
	response := []Chirp{}
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondError(w, http.StatusNotFound, "Error parsing ID", err)
		return
	}
	userChirp, err := cfg.db.GetChirpByID(r.Context(), chirpId)
	if err != nil {
		respondError(w, http.StatusNotFound, "Unable to find Chirp", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        userChirp.ID,
		CreatedAt: userChirp.CreatedAt,
		UpdatedAt: userChirp.UpdatedAt,
		Body:      userChirp.Body,
		UserID:    userChirp.UserID,
	})
}
