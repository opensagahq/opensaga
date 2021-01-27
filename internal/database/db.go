package database

import (
	"database/sql"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func Open() *sql.DB {
	const (
		dbMaxOpenConns = 4
		dbMaxIdleConns = 4
	)

	cfg, _ := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	cfg.PreferSimpleProtocol = true
	cfg.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
	}

	db := stdlib.OpenDB(*cfg)
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	return db
}
