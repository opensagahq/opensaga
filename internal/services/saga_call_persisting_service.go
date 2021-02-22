package services

import (
	"context"
	"database/sql"
	"fmt"

	"opensaga/internal/dto"
	"opensaga/internal/entities"
)

func (svc *sagaCallPersistingService) Persist(ctx context.Context, sc *dto.SagaCallCreateDTO) (err error) {
	var tx *sql.Tx

	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			if tx != nil {
				_ = tx.Rollback()
			}
		}
	}()

	tx, err = svc.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	var sagaID string
	stmt := svc.sagaIDFinder.FindIDByNameStmt(sc.Saga)
	err = tx.QueryRowContext(ctx, stmt.Query(), stmt.Args()...).Scan(&sagaID)
	if err != nil {
		return
	}

	sagaCallID := svc.uuidGenerationFunc()
	sagaCall := entities.NewSagaCall(sagaCallID, sc.IdempotencyKey, sagaID, sc.Content)

	stmt = svc.sagaCallSaver.SaveStmt(sagaCall)

	_, err = tx.ExecContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return
	}

	stmt = svc.sagaStepFinder.FindIDAndNameOfInitialStepBySagaIDStmt(sagaID)
	var sagaStepID, sagaStepName string
	err = tx.QueryRowContext(ctx, stmt.Query(), stmt.Args()...).Scan(&sagaStepID, &sagaStepName)
	if err != nil {
		return
	}

	sagaCallStep := entities.NewUnenqueuedSagaCallStep(
		sagaStepID,
		sagaCallID,
		svc.payloadExtractionFunc(sagaCall.Content, fmt.Sprintf(`step_list.%s.payload`, sagaStepName)),
	)

	stmt = svc.sagaCallStepEnqueuer.EnqueueStmt(sagaCallStep)
	_, err = tx.ExecContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return
	}

	return
}

func NewSagaCallPersistingService(cfg SagaCallPersistingServiceCfg) *sagaCallPersistingService {
	return &sagaCallPersistingService{
		sagaIDFinder:          cfg.SagaIDFinder,
		sagaCallSaver:         cfg.SagaCallSaver,
		sagaStepFinder:        cfg.SagaStepFinder,
		sagaCallStepEnqueuer:  cfg.SagaCallStepEnqueuer,
		uuidGenerationFunc:    cfg.UUIDGenerationFunc,
		payloadExtractionFunc: cfg.PayloadExtractionFunc,
		db:                    cfg.DB,
	}
}

type SagaCallPersistingServiceCfg struct {
	SagaIDFinder          SagaIDFinder
	SagaCallSaver         SagaCallSaver
	SagaStepFinder        SagaStepFinder
	SagaCallStepEnqueuer  SagaCallStepEnqueuer
	UUIDGenerationFunc    func() string
	PayloadExtractionFunc func(content, path string) string
	DB                    *sql.DB
}

type sagaCallPersistingService struct {
	sagaIDFinder          SagaIDFinder
	sagaCallSaver         SagaCallSaver
	sagaStepFinder        SagaStepFinder
	sagaCallStepEnqueuer  SagaCallStepEnqueuer
	uuidGenerationFunc    func() string
	payloadExtractionFunc func(content, path string) string
	db                    *sql.DB
}
