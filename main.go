package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type apiConfig struct {
	FileServerHits int
}

func newApiConfig() apiConfig {
	return apiConfig{FileServerHits: 0}
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := newApiConfig()
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerAdmin)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", c.FileServerHits)))
}

func (c *apiConfig) handlerAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("./admin.html")
	if err != nil {
		log.Printf("Error parsing html file: %v", err)
	}

	err = tmpl.Execute(w, c)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.FileServerHits++
		next.ServeHTTP(w, r)
	})
}
