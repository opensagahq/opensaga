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
	sagaCallRepository := repositories.NewSagaCallRepository()

	healthzHandler := healthz.New()
	sagaCreateHandler := api.NewSagaCreateHandler(api.SagaCreateHandlerCfg{
		SagaRepository:     sagaRepository,
		SagaStepRepository: sagaStepRepository,
	})
	sagaCallCreateHandler := api.NewSagaCallCreateHandler(api.SagaCallCreateHandlerCfg{
		SagaCallRepository: sagaCallRepository,
	})

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/api/saga-create", sagaCreateHandler)
	mux.Handle("/api/saga-call-create", sagaCallCreateHandler)

	_ = http.ListenAndServe(":9000", mux)
}
