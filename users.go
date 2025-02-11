package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/snowkittyselene/chirpy/internal/auth"
	"github.com/snowkittyselene/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
}

func (cfg *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userToCreate := struct {
		Email    string
		Password string
	}{}
	if err := decoder.Decode(&userToCreate); err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}
	hashedPassword, err := auth.HashPassword(userToCreate.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	u, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userToCreate.Email,
		HashedPassword: hashedPassword,
	})
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userToLogin := struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}{}
	if err := decoder.Decode(&userToLogin); err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), userToLogin.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if err = auth.CheckPasswordHash(userToLogin.Password, user.HashedPassword); err != nil {
		respondError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	expiry := 0
	if userToLogin.ExpiresInSeconds == nil {
		expiry = 3600
	} else if *userToLogin.ExpiresInSeconds > 3600 {
		expiry = 3600
	} else {
		expiry = *userToLogin.ExpiresInSeconds
	}
	expiresIn := time.Duration(expiry) * time.Second
	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error making token", err)
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
