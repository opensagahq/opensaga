package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"opensaga/internal/dto"
)

func (h *sagaCallCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	defer func() {
		if err == nil {
			SuccessResponse(w, http.StatusOK)
		} else {
			FailureResponse(w, http.StatusBadRequest, err)
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
