package repository

import (
	"context"
	"fmt"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

type WithdrawlsRepoImpl struct {
	repos *Repos
}

func NewWithdrawalsRepo(repos *Repos) *WithdrawlsRepoImpl {
	return &WithdrawlsRepoImpl{repos: repos}
}

func (r *WithdrawlsRepoImpl) Create(
	ctx context.Context,
	model *models.Withdrawal,
) (*models.Withdrawal, error) {
	// Setup transaction
	tx, err := r.repos.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to begin transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Debug(context.Background(), fmt.Sprintf("failed to rollback: %s", err))
		}
	}()

	// Making balance update
	qu, _, err := goqu.
		Update(balanceTName).
		Set(
			map[string]interface{}{
				"current":   goqu.L(fmt.Sprintf("current - %f", model.Sum)),
				"withdrawn": goqu.L(fmt.Sprintf("withdrawn + %f", model.Sum)),
			},
		).
		Where(goqu.C("user_id").Eq(model.UserID)).
		Returning(
			"balance_id",
			"user_id",
			"current",
			"withdrawn",
			"created_at",
			"updated_at",
		).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	row := tx.QueryRowxContext(ctx, qu)
	if row.Err() != nil {
		return nil, errors.Wrapf(row.Err(), "update balance error during execute query")
	}

	var balance models.Balance
	err = row.StructScan(&balance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to scan balance")
	}

	if balance.Current < 0 {
		err := tx.Rollback()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to rollback transaction")
		}

		return nil, exceptions.ErrBalanceIsNegative
	}

	// Inserting withdrawal
	model.ProcessedAt = time.Now().UTC().Format(time.RFC3339)
	qu, _, err = goqu.
		Insert(withdrawalsTName).
		Rows(model).
		Returning(
			"withdrawal_id",
			"user_id",
			"order",
			"sum",
			"processed_at",
		).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var withdrawal models.Withdrawal
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&withdrawal)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to commit")
	}

	return &withdrawal, nil
}

func (r *WithdrawlsRepoImpl) UserWithdrawals(
	ctx context.Context,
	userID string,
) (*[]models.Withdrawal, error) {
	qu, _, err := goqu.
		Select(&models.Withdrawal{}).
		From(withdrawalsTName).
		Where(goqu.C("user_id").Eq(userID)).
		Order(goqu.I("processed_at").Asc()).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	rows, err := r.repos.DB.QueryxContext(ctx, qu)
	if err != nil {
		return nil, errors.Wrapf(err, "read withdrawals error during querying")
	}

	withdrawals := []models.Withdrawal{}
	for rows.Next() {
		withdrawal := models.Withdrawal{}
		err := rows.StructScan(&withdrawal)
		if err != nil {
			return nil, errors.Wrapf(err, "read withdrawals error during scan rows")
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("read withdrawals error during querying: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(ctx, "failed to close rows", err)
		}
	}()

	return &withdrawals, nil
}
