package main

import "net/http"

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	c.FileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
