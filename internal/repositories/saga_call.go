package repositories

import (
	"opensaga/internal/entities"
)

func (r *sagaCallRepository) SaveStmt(sagaCall *entities.SagaCall) *Stmt {
	return NewStmt(
		`insert into "opensaga"."saga_call" ("idempotency_key", "saga_id", "content") values ($1, $2, $3)`,
		sagaCall.IdempotencyKey,
		sagaCall.SagaID,
		sagaCall.Content,
	)
}

func NewSagaCallRepository() *sagaCallRepository {
	return &sagaCallRepository{}
}

type sagaCallRepository struct {
}
