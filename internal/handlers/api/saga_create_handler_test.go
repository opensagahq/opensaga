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
)

func TestSagaCreateHandler_ServeHTTP(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		sagaPersistingServiceMock := new(mocks.SagaPersistingService)
		sagaPersistingServiceMock.
			On(`Persist`, context.Background(), mock.Anything).
			Return(nil)

		sut := NewSagaCreateHandler(SagaCreateHandlerCfg{
			SagaPersistingService: sagaPersistingServiceMock,
		})

		req, _ := http.NewRequest(http.MethodPost, `/api/saga-create`, bytes.NewBufferString(body))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, `unexpected status code`)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run(`expected error on persist`, func(t *testing.T) {
		sagaPersistingServiceMock := new(mocks.SagaPersistingService)
		sagaPersistingServiceMock.
			On(`Persist`, context.Background(), mock.Anything).
			Return(errors.New(`some persisting error`))

		sut := NewSagaCreateHandler(SagaCreateHandlerCfg{
			SagaPersistingService: sagaPersistingServiceMock,
		})

		req, _ := http.NewRequest(http.MethodPost, `/api/saga-create`, bytes.NewBufferString(body))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, `unexpected status code`)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run(`invalid input`, func(t *testing.T) {
		sut := NewSagaCreateHandler(SagaCreateHandlerCfg{
		})

		req, _ := http.NewRequest(http.MethodPost, `/api/saga-create`, bytes.NewBufferString(`invalid json`))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, `unexpected status code`)
	})
}

var (
	body = `{
    "id": "f5224279-d8f3-4073-bd52-2dd20b38d56b",
    "name": "saga 1",
    "step_list": [
        {
            "id": "1fc82948-5b5e-4c27-bf61-ce8cf5a66ed6",
            "next_on_success": "7a977d5c-68d4-486d-abf0-0e699fab8b18",
            "is_initial": true,
            "name": "withdraw",
            "endpoint": "https://wallet.svc.local/withdraw"
        },
        {
            "id": "7a977d5c-68d4-486d-abf0-0e699fab8b18",
            "next_on_success": "a7196b36-0c2e-4b04-9d2e-764a15e38c36",
            "next_on_failure": "eb29dae8-3e2b-40b3-aa23-d18a7d656075",
            "name": "enable paid feature",
            "endpoint": "https://paid-feature-catalog.svc.local/enable-paid-feature"
        },
        {
            "id": "a7196b36-0c2e-4b04-9d2e-764a15e38c36",
            "name": "notify",
            "endpoint": "https://notificator.svc.local/notify"
        },
        {
            "id": "eb29dae8-3e2b-40b3-aa23-d18a7d656075",
            "next_on_success": "",
            "next_on_failure": "",
            "name": "refund",
            "endpoint": "https://wallet.svc.local/refund"
        }
    ]
}
`
)
