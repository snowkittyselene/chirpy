package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	email := struct{ Email string }{}
	if err := decoder.Decode(&email); err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}
	u, err := cfg.queries.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't add user to database", err)
		return
	}
	user := User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)
}
