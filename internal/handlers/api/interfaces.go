package api

import (
	"context"

	"opensaga/internal/entities"
	"opensaga/internal/repositories"
)

type SagaRepository interface {
	SaveStmt(saga *entities.Saga) *repositories.Stmt
}

type SagaStepRepository interface {
	SaveStmt(sagaStep *entities.SagaStep) *repositories.Stmt
}

type SagaCallRepository interface {
	SaveStmt(sagaCall *entities.SagaCall) *repositories.Stmt
}

type Coordinator interface {
	Transactional(ctx context.Context, stmts ...*repositories.Stmt) (err error)
}
