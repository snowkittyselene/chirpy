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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
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
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userToLogin := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error making token", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error making refresh token", err)
		return
	}
	if _, err := cfg.db.MakeRefreshToken(r.Context(), database.MakeRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	}); err != nil {
		respondError(w, http.StatusInternalServerError, "Error adding token to database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Error getting token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := decoder.Decode(&credentials); err != nil {
		respondError(w, http.StatusUnauthorized, "Error getting request", err)
		return
	}
	hashedPassword, err := auth.HashPassword(credentials.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Error hashing password", err)
		return
	}
	newCreds, err := cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		ID:             userID,
		Email:          credentials.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Error updating credentials", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:          newCreds.ID,
		CreatedAt:   newCreds.CreatedAt,
		UpdatedAt:   newCreds.UpdatedAt,
		Email:       newCreds.Email,
		IsChirpyRed: newCreds.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	request := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = cfg.db.UpgradeUser(r.Context(), userID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
