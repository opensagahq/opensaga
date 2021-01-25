package repositories

import (
	"opensaga/internal/entities"
)

func (r *sagaRepository) SaveStmt(saga *entities.Saga) *Stmt {
	return NewStmt(`insert into "opensaga"."saga" ("id", "name") values ($1, $2)`, saga.ID, saga.Name)
}

func NewSagaRepository() *sagaRepository {
	return &sagaRepository{}
}

type sagaRepository struct {
}
