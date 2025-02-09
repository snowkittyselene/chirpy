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
