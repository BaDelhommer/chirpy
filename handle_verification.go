package main

import (
	"encoding/json"
	"net/http"
)

func decodeChirp(r *http.Request) (chirp, error) {
	var newChirp chirp

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newChirp)
	return newChirp, err
}

func sendErrorMessage(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func sendSuccessMessage(w http.ResponseWriter, resp map[string]bool) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
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

	sendSuccessMessage(w, map[string]bool{"valid": true})
}
