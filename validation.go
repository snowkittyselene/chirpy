package main

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/snowkittyselene/chirpy/internal/auth"
)

var badWords = []string{"kerfuffle", "sharbert", "fornax"}

func removeBadWords(original string) string {
	goodWords := []string{}
	words := strings.Fields(original)
	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			goodWords = append(goodWords, "****")
		} else {
			goodWords = append(goodWords, word)
		}
	}
	return strings.Join(goodWords, " ")
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Error getting token from headers", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Error retrieving user from database", err)
		return
	}
	if user.RevokedAt.Valid {
		if time.Now().Compare(user.RevokedAt.Time) > -1 {
			respondError(w, http.StatusUnauthorized, "User token revoked, cannot refresh", nil)
			return
		}
	}
	if time.Now().Compare(user.ExpiresAt) > -1 {
		respondError(w, http.StatusUnauthorized, "User token expired, cannot refresh", nil)
		return
	}
	newToken, err := auth.MakeJWT(user.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Error making new access token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: newToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken((r.Header))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error getting token from headers", err)
		return
	}
	if err = cfg.db.RevokeToken(r.Context(), token); err != nil {
		respondError(w, http.StatusInternalServerError, "Error revoking token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
