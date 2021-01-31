package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "opensaga/internal/mocks/handlers/api"
	"opensaga/internal/services"
)

func TestSagaCallCreateHandler_ServeHTTP(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		sagaCallPersistingServiceMock := new(mocks.SagaCallPersistingService)
		sagaCallPersistingServiceMock.
			On(`Persist`, context.Background(), mock.Anything).
			Return(nil)

		sut := NewSagaCallCreateHandler(SagaCallCreateHandlerCfg{
			SagaCallPersistingService: sagaCallPersistingServiceMock,
		})

		req, _ := http.NewRequest(http.MethodPost, `/api/saga-call-create`, bytes.NewBufferString(sagaCallCreateBody))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, `unexpected status code`)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run(`expected error on persist`, func(t *testing.T) {
		sagaCallPersistingServiceMock := new(mocks.SagaCallPersistingService)
		sagaCallPersistingServiceMock.
			On(`Persist`, context.Background(), mock.Anything).
			Return(errors.New(`some persisting error`))

		sut := NewSagaCallCreateHandler(SagaCallCreateHandlerCfg{
			SagaCallPersistingService: sagaCallPersistingServiceMock,
		})

		req, _ := http.NewRequest(http.MethodPost, `/api/saga-call-create`, bytes.NewBufferString(sagaCallCreateBody))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, `unexpected status code`)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run(`invalid input`, func(t *testing.T) {
		sut := NewSagaCallCreateHandler(SagaCallCreateHandlerCfg{
			SagaCallPersistingService: services.NewSagaCallPersistingService(services.SagaCallPersistingServiceCfg{}),
		})

		req, _ := http.NewRequest(http.MethodPost, `/api/saga-call-create`, bytes.NewBufferString(`invalid json`))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, `unexpected status code`)
	})
}

var (
	sagaCallCreateBody = `{
    "idempotency_key": "568521fd-0b7f-4024-ac4f-3e686e3f19e9",
    "saga": "saga 1"
}`
)
