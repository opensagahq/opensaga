package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensaga/internal/dto"
	"opensaga/internal/repositories"
)

func TestSagaPersistingService_Persist(t *testing.T) {
	t.Run(`complex`, func(t *testing.T) {
		type (
			dbMock struct {
				db        *sql.DB
				mock      sqlmock.Sqlmock
				err       error
				deferFunc func()
			}

			testCase struct {
				name     string
				expected error
				dbMock
			}
		)

		testCases := []testCase{
			{
				name:     `positive`,
				expected: nil,
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `unexpected error`)

					mock.ExpectBegin()
					mock.
						ExpectExec(`insert into "opensaga"."saga" ("id", "name") values ($1, $2)`).
						WithArgs(`2d9ed26b-3680-423a-bab1-69b80e974632`, `saga 1`).
						WillReturnResult(
							sqlmock.NewResult(0, 1),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_step"
    ("id", "saga_id", "next_on_success", "next_on_failure", "is_initial", "name", "endpoint")
    values
    ($1, $2, $3, $4, $5, $6, $7)`).
						WithArgs(
							`ad754e94-6276-4bbe-9cde-72aaaa642074`,
							`2d9ed26b-3680-423a-bab1-69b80e974632`,
							nil,
							nil,
							true,
							`saga 1 / step 1`,
							`https://foobar.svc.local`,
						).
						WillReturnResult(
							sqlmock.NewResult(0, 1),
						)
					mock.ExpectCommit()

					return dbMock{
						db:   db,
						mock: mock,
						err:  err,
						deferFunc: func(db *sql.DB) func() {
							return func() {
								db.Close()
							}
						}(db),
					}
				}(),
			},
			{
				name:     `can not start tx'`,
				expected: errors.New(`begin error`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `begin error`)

					mock.
						ExpectBegin().
						WillReturnError(
							errors.New(`begin error`),
						)

					return dbMock{
						db:   db,
						mock: mock,
						err:  err,
						deferFunc: func(db *sql.DB) func() {
							return func() {
								db.Close()
							}
						}(db),
					}
				}(),
			},
			{
				name:     `error saving saga`,
				expected: errors.New(`error saving saga`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `begin error`)

					mock.ExpectBegin()
					mock.
						ExpectExec(`insert into "opensaga"."saga" ("id", "name") values ($1, $2)`).
						WithArgs(`2d9ed26b-3680-423a-bab1-69b80e974632`, `saga 1`).
						WillReturnError(
							errors.New(`error saving saga`),
						)
					mock.ExpectRollback()

					return dbMock{
						db:   db,
						mock: mock,
						err:  err,
						deferFunc: func(db *sql.DB) func() {
							return func() {
								db.Close()
							}
						}(db),
					}
				}(),
			},
			{
				name:     `error saving saga step`,
				expected: errors.New(`error saving saga step`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `begin error`)

					mock.ExpectBegin()
					mock.
						ExpectExec(`insert into "opensaga"."saga" ("id", "name") values ($1, $2)`).
						WithArgs(`2d9ed26b-3680-423a-bab1-69b80e974632`, `saga 1`).
						WillReturnResult(
							sqlmock.NewResult(0, 1),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_step"
    ("id", "saga_id", "next_on_success", "next_on_failure", "is_initial", "name", "endpoint")
    values
    ($1, $2, $3, $4, $5, $6, $7)`).
						WithArgs(
							`ad754e94-6276-4bbe-9cde-72aaaa642074`,
							`2d9ed26b-3680-423a-bab1-69b80e974632`,
							nil,
							nil,
							true,
							`saga 1 / step 1`,
							`https://foobar.svc.local`,
						).
						WillReturnError(errors.New(`error saving saga step`))
					mock.ExpectRollback()

					return dbMock{
						db:   db,
						mock: mock,
						err:  err,
						deferFunc: func(db *sql.DB) func() {
							return func() {
								db.Close()
							}
						}(db),
					}
				}(),
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				defer testCase.deferFunc()

				ctx := context.Background()

				sut := NewSagaPersistingService(SagaPersistingServiceCfg{
					SagaSaver:     repositories.NewSagaRepository(),
					SagaStepSaver: repositories.NewSagaStepRepository(),
					DB:            testCase.db,
				})

				s := &dto.SagaCreateDTO{
					ID:   "2d9ed26b-3680-423a-bab1-69b80e974632",
					Name: "saga 1",
					StepList: []*dto.SagaStepCreateDTO{
						{
							ID:            "ad754e94-6276-4bbe-9cde-72aaaa642074",
							NextOnSuccess: nil,
							NextOnFailure: nil,
							IsInitial:     true,
							Name:          "saga 1 / step 1",
							Endpoint:      "https://foobar.svc.local",
						},
					},
				}

				err := sut.Persist(ctx, s)

				assert.Equal(t, testCase.expected, err)
				assert.NoError(t, testCase.mock.ExpectationsWereMet(), `unmet expectations`)
			})
		}
	})
}
