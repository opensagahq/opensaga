package repositories

import "opensaga/internal/entities"

func (r *sagaStepRepository) Save(sagaStep *entities.SagaStep) error {
	return nil
}

func NewSagaStepRepository() *sagaStepRepository {
	return &sagaStepRepository{}
}

type sagaStepRepository struct {
}
