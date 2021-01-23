package repositories

import (
	"context"
	"testing"
)

func TestSagaStepRepository_Save(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		ctx := context.Background()

		sut := NewSagaStepRepository()

		err := sut.Save(ctx, nil)
		if err != nil {
			t.Errorf(`unexpected error %s`, err)
		}
	})
}
