package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"gophermart/internal/closer"
	"gophermart/internal/config"
	"gophermart/internal/log"

	_ "github.com/jackc/pgx/stdlib"
)

func Migrate(ctx context.Context, cfg *config.Postgres) {
	db, err := goose.OpenDBWithDriver("pgx", cfg.DSN)
	if err != nil {
		log.Fatal(ctx, "failed to open DB for migration", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(ctx, "failed to close DB after migration", err)
		}
	}()

	if err := goose.RunContext(
		ctx,
		"up",
		db,
		cfg.MigrationFolder,
	); err != nil {
		log.Error(ctx, "failed to run up DB migration", err)
	}
}

func InitDB(ctx context.Context, cfg *config.Postgres) (*sqlx.DB, error) {
	log.Info(ctx, "migrating db")
	Migrate(ctx, cfg)

	log.Info(ctx, "initializing db")
	db, err := sqlx.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("make sql db instance error: %w", err)
	}
	log.Info(ctx, "db initialized")

	closer.Add(db.Close)

	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.IdleConn)

	go func() {
		ticker := time.NewTicker(cfg.PingInterval)
		for {
			select {
			case <-ticker.C:
				if err := db.Ping(); err != nil {
					log.Error(ctx, "error ping db", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return db, nil
}
