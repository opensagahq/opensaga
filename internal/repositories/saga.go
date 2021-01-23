package repositories

import "opensaga/internal/entities"

func (r *sagaRepository) Save(saga *entities.Saga) error {
	return nil
}

func NewSagaRepository() *sagaRepository {
	return &sagaRepository{}
}

type sagaRepository struct {
}
