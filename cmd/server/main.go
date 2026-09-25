package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/newyear-days/internal/api"
)

func main() {
	addr := listenAddr()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/days-until-new-year", api.DaysUntilNewYearHandler)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}
