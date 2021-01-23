package api

import (
	"encoding/json"
	"net/http"

	"opensaga/internal/dto"
	"opensaga/internal/entities"
)

func (h *sagaCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sagaDTO dto.SagaCreateDTO

	err := json.NewDecoder(r.Body).Decode(&sagaDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "failed"}`))

		return
	}

	// todo extract into SagaPersister
	// todo tx
	saga := &entities.Saga{
		ID:   sagaDTO.ID,
		Name: sagaDTO.Name,
	}
	_ = h.sagaRepository.Save(saga)

	for _, sagaStepDTO := range sagaDTO.StepList {
		sagaStep := &entities.SagaStep{
			ID:            sagaStepDTO.ID,
			SagaID:        sagaDTO.ID,
			NextOnSuccess: sagaStepDTO.NextOnSuccess,
			NextOnFailure: sagaStepDTO.NextOnFailure,
			IsInitial:     sagaStepDTO.IsInitial,
			Name:          sagaStepDTO.Name,
			Endpoint:      sagaStepDTO.Endpoint,
		}

		_ = h.sagaStepRepository.Save(sagaStep)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func NewSagaCreateHandler(cfg SagaCreateHandlerCfg) *sagaCreateHandler {
	return &sagaCreateHandler{
		sagaRepository:     cfg.SagaRepository,
		sagaStepRepository: cfg.SagaStepRepository,
	}
}

type SagaCreateHandlerCfg struct {
	SagaRepository     SagaRepository
	SagaStepRepository SagaStepRepository
}

type sagaCreateHandler struct {
	sagaRepository     SagaRepository
	sagaStepRepository SagaStepRepository
}

type SagaRepository interface {
	Save(saga *entities.Saga) error
}

type SagaStepRepository interface {
	Save(sagaStep *entities.SagaStep) error
}
