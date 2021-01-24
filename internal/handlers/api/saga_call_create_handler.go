package api

import (
	"context"
	"encoding/json"
	"net/http"

	"opensaga/internal/dto"
	"opensaga/internal/entities"
)

func (h *sagaCallCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sagaCallCreateDTO dto.SagaCallCreateDTO

	err := json.NewDecoder(r.Body).Decode(&sagaCallCreateDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "failed"}`))

		return
	}

	sagaCall := &entities.SagaCall{
		IdempotencyKey: sagaCallCreateDTO.IdempotencyKey,
		SagaID:         sagaCallCreateDTO.Saga,
	}
	_ = h.sagaCallRepository.Save(ctx, sagaCall)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func NewSagaCallCreateHandler(cfg SagaCallCreateHandlerCfg) *sagaCallCreateHandler {
	return &sagaCallCreateHandler{
		sagaCallRepository: cfg.SagaCallRepository,
	}
}

type SagaCallCreateHandlerCfg struct {
	SagaCallRepository SagaCallRepository
}

type sagaCallCreateHandler struct {
	sagaCallRepository SagaCallRepository
}

type SagaCallRepository interface {
	Save(ctx context.Context, sagaCall *entities.SagaCall) error
}
