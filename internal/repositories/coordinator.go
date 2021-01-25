package repositories

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func (c *coordinator) Transactional(ctx context.Context, stmts ...*Stmt) (err error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit(ctx)
		} else {
			_ = tx.Rollback(ctx)
		}
	}()

	for _, stmt := range stmts {
		_, err := tx.Exec(ctx, stmt.Query(), stmt.Args()...)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewCoordinator(cfg CoordinatorCfg) *coordinator {
	return &coordinator{
		conn: cfg.Conn,
	}
}

type CoordinatorCfg struct {
	Conn *pgxpool.Pool
}

type coordinator struct {
	conn *pgxpool.Pool
}
