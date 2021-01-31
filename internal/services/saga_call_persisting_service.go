package services

import (
	"context"
	"database/sql"

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

	sagaCall := entities.NewSagaCall(sc.IdempotencyKey, sagaID, sc.Content)

	stmt = svc.sagaCallSaver.SaveStmt(sagaCall)

	_, err = tx.ExecContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return
	}

	return
}

func NewSagaCallPersistingService(cfg SagaCallPersistingServiceCfg) *sagaCallPersistingService {
	return &sagaCallPersistingService{
		sagaIDFinder:  cfg.SagaIDFinder,
		sagaCallSaver: cfg.SagaCallSaver,
		db:            cfg.DB,
	}
}

type SagaCallPersistingServiceCfg struct {
	SagaIDFinder  SagaIDFinder
	SagaCallSaver SagaCallSaver
	DB            *sql.DB
}

type sagaCallPersistingService struct {
	sagaIDFinder  SagaIDFinder
	sagaCallSaver SagaCallSaver
	db            *sql.DB
}
