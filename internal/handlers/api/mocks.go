// todo replace handcrafted mocks with generated
package api

import (
	"context"

	"opensaga/internal/repositories"
)

func (c *coordinatorMock) Transactional(ctx context.Context, stmts ...*repositories.Stmt) (err error) {
	return nil
}

func NewCoordinatorMock() *coordinatorMock {
	return &coordinatorMock{}
}

type coordinatorMock struct {
}
