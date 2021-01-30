package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"opensaga/internal/dto"
)

func (h *sagaCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	defer func() {
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status": "failed", "error": "%s"}`, err)))
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
