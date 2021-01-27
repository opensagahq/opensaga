package repositories

import (
	"context"
	"database/sql"
)

func (c *coordinator) Transactional(ctx context.Context, stmts ...*Stmt) (err error) {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	for _, stmt := range stmts {
		_, err := tx.ExecContext(ctx, stmt.Query(), stmt.Args()...)
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
	Conn *sql.DB
}

type coordinator struct {
	conn *sql.DB
}
