package repositories

import (
	"opensaga/internal/entities"
)

func (r *sagaCallStepRepository) EnqueueStmt(sagaCallStep *entities.SagaCallStep) *Stmt {
	return NewStmt(`insert into "opensaga"."saga_call_step_queue" ("saga_step_id", "saga_call_id", "payload") values ($1, $2, $3)`, sagaCallStep.SagaStepID, sagaCallStep.SagaCallID, sagaCallStep.Payload)
}

func NewSagaCallStepRepository() *sagaCallStepRepository {
	return &sagaCallStepRepository{}
}

type sagaCallStepRepository struct {
}
