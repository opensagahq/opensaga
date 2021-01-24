package repositories

import (
	"context"

	"opensaga/internal/entities"
)

func (r *sagaCallRepository) Save(ctx context.Context, sagaCall *entities.SagaCall) error {
	return nil
}

func NewSagaCallRepository() *sagaCallRepository {
	return &sagaCallRepository{}
}

type sagaCallRepository struct {
}
