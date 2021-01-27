package main

import (
	"net/http"

	"opensaga/internal/database"
	"opensaga/internal/handlers/api"
	"opensaga/internal/handlers/healthz"
	"opensaga/internal/repositories"
)

func main() {
	db := database.Open()

	coordinator := repositories.NewCoordinator(repositories.CoordinatorCfg{Conn: db})
	sagaRepository := repositories.NewSagaRepository()
	sagaStepRepository := repositories.NewSagaStepRepository()
	sagaCallRepository := repositories.NewSagaCallRepository()

	healthzHandler := healthz.New()
	sagaCreateHandler := api.NewSagaCreateHandler(api.SagaCreateHandlerCfg{
		SagaRepository:     sagaRepository,
		SagaStepRepository: sagaStepRepository,
		Coordinator:        coordinator,
	})
	sagaCallCreateHandler := api.NewSagaCallCreateHandler(api.SagaCallCreateHandlerCfg{
		SagaCallRepository: sagaCallRepository,
		Coordinator:        coordinator,
	})

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/api/saga-create", sagaCreateHandler)
	mux.Handle("/api/saga-call-create", sagaCallCreateHandler)

	_ = http.ListenAndServe(":9000", mux)
}
