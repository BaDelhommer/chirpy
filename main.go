package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func newApiConfig() apiConfig {
	return apiConfig{fileServerHits: 0}
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := newApiConfig()
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", c.fileServerHits)))
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits++
		next.ServeHTTP(w, r)
	})
}
