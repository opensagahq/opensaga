package repositories

import (
	"context"
	"testing"
)

func TestSagaRepository_Save(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		ctx := context.Background()

		sut := NewSagaRepository()

		err := sut.SaveStmt(ctx, nil)
		if err != nil {
			t.Errorf(`unexpected error %s`, err)
		}
	})
}
