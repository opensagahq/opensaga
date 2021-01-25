package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"opensaga/internal/handlers/api"
	"opensaga/internal/handlers/healthz"
	"opensaga/internal/repositories"
)

func main() {
	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to parse DATABASE_URL: %v", err)
	}
	poolConfig.ConnConfig.PreferSimpleProtocol = true

	db, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

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
