package api

import (
	"context"

	"opensaga/internal/dto"
)

type SagaPersistingService interface {
	Persist(context.Context, *dto.SagaCreateDTO) (err error)
}

type SagaCallPersistingService interface {
	Persist(context.Context, *dto.SagaCallCreateDTO) (err error)
}
