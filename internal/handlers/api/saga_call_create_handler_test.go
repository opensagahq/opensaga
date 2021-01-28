package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"opensaga/internal/repositories"
)

func TestSagaCallCreateHandler_ServeHTTP(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		sut := NewSagaCallCreateHandler(SagaCallCreateHandlerCfg{
			SagaCallRepository: repositories.NewSagaCallRepository(),
			Coordinator:        NewCoordinatorMock(),
		})
		req, _ := http.NewRequest(http.MethodPost, `/api/saga-call-create`, bytes.NewBufferString(sagaCallCreateBody))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf(`unexpected status code: got "%v", want "%v"`, status, http.StatusOK)
		}
	})

	t.Run(`invalid json input`, func(t *testing.T) {
		sut := NewSagaCallCreateHandler(SagaCallCreateHandlerCfg{
			SagaCallRepository: repositories.NewSagaCallRepository(),
			Coordinator:        NewCoordinatorMock(),
		})
		req, _ := http.NewRequest(http.MethodPost, `/api/saga-call-create`, bytes.NewBufferString(`invalid json`))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf(`unexpected status code: got "%v", want "%v"`, status, http.StatusBadRequest)
		}
	})
}

var (
	sagaCallCreateBody = `{
    "idempotency_key": "568521fd-0b7f-4024-ac4f-3e686e3f19e9",
    "saga": "saga 1"
}`
)
