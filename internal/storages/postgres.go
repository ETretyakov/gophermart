package storages

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"gophermart/internal/log"
	"io"
	"time"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

var recoverableErrors = []error{
	pgx.ErrDeadConn,
	sql.ErrConnDone,
	pgx.ErrConnBusy,
	driver.ErrBadConn,
	io.EOF,
}

type Retryer struct {
	Strategy Strategy
	OnRetry  func(ctx context.Context, n int, err error)
}

type SQLDB struct {
	db      *sqlx.DB
	retryer *Retryer
}

func NewSQLDB(db *sqlx.DB) *SQLDB {
	if db == nil {
		return nil
	}

	r := &SQLDB{
		db: db,
	}

	r.retryer = &Retryer{
		Strategy: Backoff(3, 5*time.Second, .5, true),
		OnRetry: func(ctx context.Context, n int, err error) {
			log.Info(ctx, fmt.Sprintf("reconnecting DB (%d): %s", n, err))
			if err = r.db.PingContext(ctx); err != nil {
				log.Error(ctx, "failed to connect to DB", err)
			}
		},
	}

	return r
}

func (r *SQLDB) retry(ctx context.Context, fn func() error) error {
	return r.retryer.Do(ctx, fn, recoverableErrors...)
}

func (r *SQLDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (tx *sqlx.Tx, err error) {
	err = r.retry(ctx, func() error {
		var err error
		tx, err = r.db.BeginTxx(ctx, opts)
		return fmt.Errorf("failed to begin transaction: %w", err)
	})
	return tx, err
}

func (r *SQLDB) QueryxContext(
	ctx context.Context,
	query string,
	args ...any,
) (*sqlx.Rows, error) {
	var rows *sqlx.Rows
	err := r.retry(ctx, func() error {
		var err error
		rows, err = r.db.QueryxContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to query context: %w", err)
		}
		err = rows.Err()
		if err != nil {
			return fmt.Errorf("failed to query context error in rows: %w", err)
		}
		return nil
	})
	return rows, err
}

func (r *SQLDB) QueryRowxContext(
	ctx context.Context,
	query string,
	args ...any,
) (row *sqlx.Row) {
	_ = r.retry(ctx, func() error {
		row = r.db.QueryRowxContext(ctx, query, args...)
		return row.Err()
	})
	return
}

func (r *SQLDB) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	var res sql.Result
	err := r.retry(ctx, func() error {
		var err error
		res, err = r.db.ExecContext(ctx, query, args...)
		return err
	})
	return res, err
}

func (r *SQLDB) PingDB(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
