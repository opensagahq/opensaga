package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"opensaga/internal/dto"
	"opensaga/internal/entities"
)

func (h *sagaCallCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sagaCallCreateDTO dto.SagaCallCreateDTO

	body, _ := ioutil.ReadAll(r.Body)

	err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&sagaCallCreateDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "failed"}`))

		return
	}

	sagaCall := &entities.SagaCall{
		IdempotencyKey: sagaCallCreateDTO.IdempotencyKey,
		SagaID:         "f5224279-d8f3-4073-bd52-2dd20b38d56b", // todo find saga by name
	}
	sagaCall.Content = string(body)

	stmt := h.sagaCallRepository.SaveStmt(sagaCall)
	err = h.coordinator.Transactional(ctx, stmt)
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

func NewSagaCallCreateHandler(cfg SagaCallCreateHandlerCfg) *sagaCallCreateHandler {
	return &sagaCallCreateHandler{
		sagaCallRepository: cfg.SagaCallRepository,
		coordinator:        cfg.Coordinator,
	}
}

type SagaCallCreateHandlerCfg struct {
	SagaCallRepository SagaCallRepository
	Coordinator        Coordinator
}

type sagaCallCreateHandler struct {
	sagaCallRepository SagaCallRepository
	coordinator        Coordinator
}
