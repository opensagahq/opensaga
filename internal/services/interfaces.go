package services

import (
	. "opensaga/internal/entities"
	. "opensaga/internal/repositories"
)

type SagaSaver interface {
	SaveStmt(*Saga) *Stmt
}

type SagaStepSaver interface {
	SaveStmt(step *SagaStep) *Stmt
}

type SagaIDFinder interface {
	FindIDByNameStmt(string) *Stmt
}

type SagaCallSaver interface {
	SaveStmt(*SagaCall) *Stmt
}
