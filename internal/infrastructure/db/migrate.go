package db

import (
	"Fridger/internal/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(pool *pgxpool.Pool, cfg configuration.MigrationsConfiguration) error {
	db := stdlib.OpenDBFromPool(pool)
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, cfg.MigrationsFolder); err != nil {
		return err
	}

	return nil
}
