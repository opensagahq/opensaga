package main

import (
	"net/http"

	"opensaga/internal/database"
	"opensaga/internal/handlers/api"
	"opensaga/internal/handlers/healthz"
	"opensaga/internal/repositories"
	"opensaga/internal/services"
)

func main() {
	db := database.Open()

	sagaRepository := repositories.NewSagaRepository()
	sagaStepRepository := repositories.NewSagaStepRepository()
	sagaCallRepository := repositories.NewSagaCallRepository()

	sagaPersistingService := services.NewSagaPersistingService(services.SagaPersistingServiceCfg{
		SagaSaver:     sagaRepository,
		SagaStepSaver: sagaStepRepository,
		DB:            db,
	})
	sagaCallPersistingService := services.NewSagaCallPersistingService(services.SagaCallPersistingServiceCfg{
		SagaIDFinder:  sagaRepository,
		SagaCallSaver: sagaCallRepository,
		DB:            db,
	})

	healthzHandler := healthz.New()
	sagaCreateHandler := api.NewSagaCreateHandler(api.SagaCreateHandlerCfg{
		SagaPersistingService: sagaPersistingService,
	})
	sagaCallCreateHandler := api.NewSagaCallCreateHandler(api.SagaCallCreateHandlerCfg{
		SagaCallPersistingService: sagaCallPersistingService,
	})

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/api/saga-create", sagaCreateHandler)
	mux.Handle("/api/saga-call-create", sagaCallCreateHandler)

	_ = http.ListenAndServe(":9000", mux)
}
