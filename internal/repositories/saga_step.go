package repositories

import (
	"context"
	"opensaga/internal/entities"
)

func (r *sagaStepRepository) Save(ctx context.Context, sagaStep *entities.SagaStep) error {
	return nil
}

func NewSagaStepRepository() *sagaStepRepository {
	return &sagaStepRepository{}
}

type sagaStepRepository struct {
}
