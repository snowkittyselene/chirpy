package main

import (
	"slices"
	"strings"
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
