package api

import (
	"encoding/json"
	"net/http"

	"opensaga/internal/dto"
)

func (h *sagaCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	defer func() {
		if err == nil {
			SuccessResponse(w, http.StatusOK)
		} else {
			FailureResponse(w, http.StatusBadRequest, err)
		}
	}()

	var sagaCreateDTO dto.SagaCreateDTO
	err = json.NewDecoder(r.Body).Decode(&sagaCreateDTO)
	if err != nil {
		return
	}

	err = h.sagaPersistingService.Persist(ctx, &sagaCreateDTO)
	if err != nil {
		return
	}
}

func NewSagaCreateHandler(cfg SagaCreateHandlerCfg) *sagaCreateHandler {
	return &sagaCreateHandler{
		sagaPersistingService: cfg.SagaPersistingService,
	}
}

type SagaCreateHandlerCfg struct {
	SagaPersistingService SagaPersistingService
}

type sagaCreateHandler struct {
	sagaPersistingService SagaPersistingService
}
