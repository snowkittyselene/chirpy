package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
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
			Valid: true,
		}
		respondWithJSON(w, http.StatusOK, successMessage)
	}

}
