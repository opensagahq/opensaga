package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCoordinator_Transactional(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		sut := NewCoordinator(CoordinatorCfg{Conn: db})

		mock.ExpectBegin()
		mock.ExpectCommit()

		_ = sut.Transactional(context.Background())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
