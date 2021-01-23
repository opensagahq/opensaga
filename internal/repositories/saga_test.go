package repositories

import "testing"

func TestSagaRepository_Save(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		sut := NewSagaRepository()

		err := sut.Save(nil)
		if err != nil {
			t.Errorf(`unexpected error %s`, err)
		}
	})
}
