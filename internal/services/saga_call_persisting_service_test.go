package services

import (
	"context"
	"testing"

	"opensaga/internal/dto"
	"opensaga/internal/repositories"
)

func TestSagaCallPersistingService_Persist(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		ctx := context.Background()

		sut := NewSagaCallPersistingService(SagaCallPersistingServiceCfg{
			// todo replace with mock
			SagaIDFinder:  repositories.NewSagaRepository(),
			SagaCallSaver: repositories.NewSagaCallRepository(),
		})

		sc := &dto.SagaCallCreateDTO{
			IdempotencyKey: "df703e35-f71b-4dea-8fe8-8da3ecc4357f",
			Saga:           "saga 1",
			Content:        "{}",
		}

		err := sut.Persist(ctx, sc)
		if err != nil {
			t.Errorf(`unexpected error "%v"`, err)
		}
	})
}
