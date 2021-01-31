package services

import (
	"context"
	"database/sql"

	"opensaga/internal/dto"
	"opensaga/internal/entities"
)

func (svc *sagaPersistingService) Persist(ctx context.Context, sg *dto.SagaCreateDTO) (err error) {
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

	saga := entities.NewSaga(sg.ID, sg.Name)

	stmt := svc.sagaSaver.SaveStmt(saga)
	_, err = tx.ExecContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return
	}

	for _, sagaStepDTO := range sg.StepList {
		sagaStep := &entities.SagaStep{
			ID:            sagaStepDTO.ID,
			SagaID:        saga.ID,
			NextOnSuccess: sagaStepDTO.NextOnSuccess,
			NextOnFailure: sagaStepDTO.NextOnFailure,
			IsInitial:     sagaStepDTO.IsInitial,
			Name:          sagaStepDTO.Name,
			Endpoint:      sagaStepDTO.Endpoint,
		}

		stmt := svc.sagaStepSaver.SaveStmt(sagaStep)
		_, err = tx.ExecContext(ctx, stmt.Query(), stmt.Args()...)
		if err != nil {
			return
		}
	}

	return
}

func NewSagaPersistingService(cfg SagaPersistingServiceCfg) *sagaPersistingService {
	return &sagaPersistingService{
		sagaSaver:     cfg.SagaSaver,
		sagaStepSaver: cfg.SagaStepSaver,
		db:            cfg.DB,
	}
}

type SagaPersistingServiceCfg struct {
	SagaSaver     SagaSaver
	SagaStepSaver SagaStepSaver
	DB            *sql.DB
}

type sagaPersistingService struct {
	sagaSaver     SagaSaver
	sagaStepSaver SagaStepSaver
	db            *sql.DB
}
