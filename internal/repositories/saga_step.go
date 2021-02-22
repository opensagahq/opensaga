package repositories

import (
	"opensaga/internal/entities"
)

func (r *sagaStepRepository) SaveStmt(sagaStep *entities.SagaStep) *Stmt {
	return NewStmt(`insert into "opensaga"."saga_step"
    ("id", "saga_id", "next_on_success", "next_on_failure", "is_initial", "name", "endpoint")
    values
    ($1, $2, $3, $4, $5, $6, $7)`,
		sagaStep.ID,
		sagaStep.SagaID,
		sagaStep.NextOnSuccess,
		sagaStep.NextOnFailure,
		sagaStep.IsInitial,
		sagaStep.Name,
		sagaStep.Endpoint,
	)
}

func (r *sagaStepRepository) FindIDAndNameOfInitialStepBySagaIDStmt(sagaID string) *Stmt {
	return NewStmt(`select "id", "name" from "opensaga"."saga_step" where "saga_id" = $1 and "is_initial"`, sagaID)
}

func NewSagaStepRepository() *sagaStepRepository {
	return &sagaStepRepository{}
}

type sagaStepRepository struct {
}
