package repositories

import "testing"

func TestSagaStepRepository_Save(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		sut := NewSagaStepRepository()

		err := sut.Save(nil)
		if err != nil {
			t.Errorf(`unexpected error %s`, err)
		}
	})
}
