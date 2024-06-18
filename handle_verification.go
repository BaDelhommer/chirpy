package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
)

func decodeChirp(r *http.Request) (database.Chirp, error) {
	var newChirp database.Chirp

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newChirp)
	return newChirp, err
}

func sendErrorMessage(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func sendSuccessMessage(w http.ResponseWriter, resp map[string]string) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func cleanChirp(newChirp database.Chirp) database.Chirp {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	lowerChirp := strings.ToLower(newChirp.Body)
	splitChirp := strings.Split(lowerChirp, " ")
	returnChirp := strings.Split(newChirp.Body, " ")

	for i := range splitChirp {
		for j := range badWords {
			if splitChirp[i] == badWords[j] {
				returnChirp[i] = "****"
			}
		}
	}

	cleanChirp := database.Chirp{
		Body: strings.Join(returnChirp, " "),
	}

	return cleanChirp
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	newChirp, err := decodeChirp(r)
	if err != nil {
		sendErrorMessage(w, "invalid request", 400)
		return
	}

	if len(newChirp.Body) > 140 {
		sendErrorMessage(w, "Chirp is too long", 400)
		return
	}

	cleanChirp := cleanChirp(newChirp)

	sendSuccessMessage(w, map[string]string{"cleaned_body": cleanChirp.Body})
}
