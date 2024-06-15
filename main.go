package main

import (
	"fmt"
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}
	fmt.Println("Server Running")
	server.ListenAndServe()
}
