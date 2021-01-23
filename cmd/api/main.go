package main

import (
	"net/http"

	"opensaga/internal/handlers/healthz"
)

func main() {
	healthzHandler := healthz.New()

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthzHandler)

	_ = http.ListenAndServe(":9000", mux)
}
