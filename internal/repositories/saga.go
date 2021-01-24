package repositories

import (
	"context"

	"opensaga/internal/entities"
)

func (r *sagaRepository) Save(ctx context.Context, saga *entities.Saga) error {
	return nil
}

func NewSagaRepository() *sagaRepository {
	return &sagaRepository{}
}

type sagaRepository struct {
}
