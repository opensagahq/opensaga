package main

import (
	"net/http"

	"opensaga/internal/handlers/api"
	"opensaga/internal/handlers/healthz"
	"opensaga/internal/repositories"
)

func main() {
	sagaRepository := repositories.NewSagaRepository()
	sagaStepRepository := repositories.NewSagaStepRepository()

	healthzHandler := healthz.New()
	sagaCreateHandler := api.NewSagaCreateHandler(api.SagaCreateHandlerCfg{
		SagaRepository:     sagaRepository,
		SagaStepRepository: sagaStepRepository,
	})

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/api/saga-create", sagaCreateHandler)

	_ = http.ListenAndServe(":9000", mux)
}
