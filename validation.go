package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

var badWords = []string{"kerfuffle", "sharbert", "fornax"}

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	sentChirp := chirp{}
	if err := decoder.Decode(&sentChirp); err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't decode Chirp", err)
		return
	}
	if len(sentChirp.Body) > 140 {
		respondError(w, http.StatusBadRequest, "Chirp is too long", nil)
	} else {
		successMessage := validResponse{
			CleanedBody: removeBadWords(sentChirp.Body),
		}
		respondWithJSON(w, http.StatusOK, successMessage)
	}

}

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
