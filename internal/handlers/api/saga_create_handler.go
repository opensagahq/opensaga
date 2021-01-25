package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"opensaga/internal/dto"
	"opensaga/internal/entities"
	"opensaga/internal/repositories"
)

func (h *sagaCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sagaDTO dto.SagaCreateDTO

	err := json.NewDecoder(r.Body).Decode(&sagaDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "failed"}`))

		return
	}

	// todo extract into SagaPersister
	stmts := make([]*repositories.Stmt, 1+len(sagaDTO.StepList))

	saga := &entities.Saga{
		ID:   sagaDTO.ID,
		Name: sagaDTO.Name,
	}
	stmts[0] = h.sagaRepository.SaveStmt(saga)

	for i, sagaStepDTO := range sagaDTO.StepList {
		sagaStep := &entities.SagaStep{
			ID:            sagaStepDTO.ID,
			SagaID:        sagaDTO.ID,
			NextOnSuccess: sagaStepDTO.NextOnSuccess,
			NextOnFailure: sagaStepDTO.NextOnFailure,
			IsInitial:     sagaStepDTO.IsInitial,
			Name:          sagaStepDTO.Name,
			Endpoint:      sagaStepDTO.Endpoint,
		}

		stmts[i+1] = h.sagaStepRepository.SaveStmt(sagaStep)
	}

	err = h.coordinator.Transactional(ctx, stmts...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status": "failed", "error": "%s"}`, err)))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func NewSagaCreateHandler(cfg SagaCreateHandlerCfg) *sagaCreateHandler {
	return &sagaCreateHandler{
		sagaRepository:     cfg.SagaRepository,
		sagaStepRepository: cfg.SagaStepRepository,
		coordinator:        cfg.Coordinator,
	}
}

type SagaCreateHandlerCfg struct {
	SagaRepository     SagaRepository
	SagaStepRepository SagaStepRepository
	Coordinator        Coordinator
}

type sagaCreateHandler struct {
	sagaRepository     SagaRepository
	sagaStepRepository SagaStepRepository
	coordinator        Coordinator
}
