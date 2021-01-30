package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"opensaga/internal/dto"
)

func (h *sagaCallCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var sagaCallCreateDTO dto.SagaCallCreateDTO
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&sagaCallCreateDTO)
	if err != nil {
		return
	}
	sagaCallCreateDTO.Content = string(body)

	err = h.sagaCallPersistingService.Persist(ctx, &sagaCallCreateDTO)
	if err != nil {
		return
	}
}

func NewSagaCallCreateHandler(cfg SagaCallCreateHandlerCfg) *sagaCallCreateHandler {
	return &sagaCallCreateHandler{
		sagaCallPersistingService: cfg.SagaCallPersistingService,
	}
}

type SagaCallCreateHandlerCfg struct {
	SagaCallPersistingService SagaCallPersistingService
}

type sagaCallCreateHandler struct {
	sagaCallPersistingService SagaCallPersistingService
}
