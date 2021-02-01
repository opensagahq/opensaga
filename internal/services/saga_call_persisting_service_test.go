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

func TestSagaCallPersistingService_Persist(t *testing.T) {
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
						ExpectQuery(`select "id" from "opensaga"."saga" where "name" = $1`).
						WithArgs(`saga 1`).
						WillReturnRows(
							sqlmock.
								NewRows([]string{`id`}).
								AddRow(`4146bc8e-4936-4771-a474-5c5824369dd6`),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_call" ("id", "idempotency_key", "saga_id", "content") values ($1, $2, $3, $4)`).
						WithArgs(`7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`, `df703e35-f71b-4dea-8fe8-8da3ecc4357f`, `4146bc8e-4936-4771-a474-5c5824369dd6`, `{}`).
						WillReturnResult(
							sqlmock.NewResult(0, 1),
						)
					mock.
						ExpectQuery(`select "id", "name" from "opensaga"."saga_step" where "saga_id" = $1 and "is_initial"`).
						WithArgs(`4146bc8e-4936-4771-a474-5c5824369dd6`).
						WillReturnRows(
							sqlmock.
								NewRows([]string{`id`, `name`}).
								AddRow(`bc22e32b-df92-4f8b-a6ea-0a1fa9f28ef9`, `saga 1`),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_call_step_queue" ("saga_step_id", "saga_call_id", "payload") values ($1, $2, $3)`).
						WithArgs(`bc22e32b-df92-4f8b-a6ea-0a1fa9f28ef9`, `7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`, `{}`).
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

					mock.ExpectBegin().WillReturnError(errors.New(`begin error`))

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
				name:     `saga not found`,
				expected: errors.New(`saga not found`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `unexpected error`)

					mock.ExpectBegin()
					mock.
						ExpectQuery(`select "id" from "opensaga"."saga" where "name" = $1`).
						WithArgs(`saga 1`).
						WillReturnError(errors.New(`saga not found`))
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
				name:     `insert error`,
				expected: errors.New(`insert error`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `unexpected error`)

					mock.ExpectBegin()
					mock.
						ExpectQuery(`select "id" from "opensaga"."saga" where "name" = $1`).
						WithArgs(`saga 1`).
						WillReturnRows(
							sqlmock.
								NewRows([]string{`id`}).
								AddRow(`4146bc8e-4936-4771-a474-5c5824369dd6`),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_call" ("id", "idempotency_key", "saga_id", "content") values ($1, $2, $3, $4)`).
						WithArgs(`7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`, `df703e35-f71b-4dea-8fe8-8da3ecc4357f`, `4146bc8e-4936-4771-a474-5c5824369dd6`, `{}`).
						WillReturnError(errors.New(`insert error`))
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
				name:     `failed to find initial step`,
				expected: errors.New(`initial step not found`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `unexpected error`)

					mock.ExpectBegin()
					mock.
						ExpectQuery(`select "id" from "opensaga"."saga" where "name" = $1`).
						WithArgs(`saga 1`).
						WillReturnRows(
							sqlmock.
								NewRows([]string{`id`}).
								AddRow(`4146bc8e-4936-4771-a474-5c5824369dd6`),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_call" ("id", "idempotency_key", "saga_id", "content") values ($1, $2, $3, $4)`).
						WithArgs(`7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`, `df703e35-f71b-4dea-8fe8-8da3ecc4357f`, `4146bc8e-4936-4771-a474-5c5824369dd6`, `{}`).
						WillReturnResult(
							sqlmock.NewResult(0, 1),
						)
					mock.
						ExpectQuery(`select "id", "name" from "opensaga"."saga_step" where "saga_id" = $1 and "is_initial"`).
						WithArgs(`4146bc8e-4936-4771-a474-5c5824369dd6`).
						WillReturnError(errors.New(`initial step not found`))
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
				name:     `positive`,
				expected: errors.New(`enqueue failed`),
				dbMock: func() dbMock {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					require.NoError(t, err, `unexpected error`)

					mock.ExpectBegin()
					mock.
						ExpectQuery(`select "id" from "opensaga"."saga" where "name" = $1`).
						WithArgs(`saga 1`).
						WillReturnRows(
							sqlmock.
								NewRows([]string{`id`}).
								AddRow(`4146bc8e-4936-4771-a474-5c5824369dd6`),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_call" ("id", "idempotency_key", "saga_id", "content") values ($1, $2, $3, $4)`).
						WithArgs(`7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`, `df703e35-f71b-4dea-8fe8-8da3ecc4357f`, `4146bc8e-4936-4771-a474-5c5824369dd6`, `{}`).
						WillReturnResult(
							sqlmock.NewResult(0, 1),
						)
					mock.
						ExpectQuery(`select "id", "name" from "opensaga"."saga_step" where "saga_id" = $1 and "is_initial"`).
						WithArgs(`4146bc8e-4936-4771-a474-5c5824369dd6`).
						WillReturnRows(
							sqlmock.
								NewRows([]string{`id`, `name`}).
								AddRow(`bc22e32b-df92-4f8b-a6ea-0a1fa9f28ef9`, `saga 1`),
						)
					mock.
						ExpectExec(`insert into "opensaga"."saga_call_step_queue" ("saga_step_id", "saga_call_id", "payload") values ($1, $2, $3)`).
						WithArgs(`bc22e32b-df92-4f8b-a6ea-0a1fa9f28ef9`, `7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`, `{}`).
						WillReturnError(errors.New(`enqueue failed`))
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

				sut := NewSagaCallPersistingService(SagaCallPersistingServiceCfg{
					SagaIDFinder:          repositories.NewSagaRepository(),
					SagaCallSaver:         repositories.NewSagaCallRepository(),
					SagaStepFinder:        repositories.NewSagaStepRepository(),
					SagaCallStepEnqueuer:  repositories.NewSagaCallStepRepository(),
					UUIDGenerationFunc:    uuidGenerationFunc,
					PayloadExtractionFunc: payloadExtractionFunc,
					DB:                    testCase.db,
				})

				sc := &dto.SagaCallCreateDTO{
					IdempotencyKey: "df703e35-f71b-4dea-8fe8-8da3ecc4357f",
					Saga:           "saga 1",
					Content:        "{}",
				}

				err := sut.Persist(ctx, sc)

				assert.Equal(t, testCase.expected, err)
				assert.NoError(t, testCase.mock.ExpectationsWereMet(), `unmet expectations`)
			})
		}
	})
}

func uuidGenerationFunc() string {
	return `7cb4d6b3-d7e6-47e1-a394-6b7e8d023edc`
}

func payloadExtractionFunc(content, path string) string {
	return `{}`
}
