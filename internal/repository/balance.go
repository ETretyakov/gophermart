package repository

import (
	"context"
	"database/sql"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

type BalanceRepoImpl struct {
	repos *Repos
}

func NewBalanceRepo(repos *Repos) *BalanceRepoImpl {
	return &BalanceRepoImpl{repos: repos}
}

func (r *BalanceRepoImpl) GetOrCreateForUser(
	ctx context.Context,
	userID string,
) (*models.Balance, error) {
	qu, _, err := goqu.
		Select(&models.Balance{}).
		From(balanceTName).
		Where(goqu.C("user_id").Eq(userID)).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var balance models.Balance
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			balance, err := r.Create(ctx, models.NewBalance(userID))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create balance")
			}

			return balance, nil
		} else {
			return nil, errors.Wrapf(err, "failed to execute query")
		}
	}

	return &balance, nil
}

func (r *BalanceRepoImpl) Get(
	ctx context.Context,
	balanceID string,
) (*models.Balance, error) {
	qu, _, err := goqu.
		Select(&models.Balance{}).
		From(balanceTName).
		Where(goqu.C("balance_id").Eq(balanceID)).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var balance models.Balance
	err = r.repos.DB.QueryRowxContext(ctx, qu).Scan(&balance)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, exceptions.ErrBalanceNotFound
	}

	return &balance, nil
}

func (r *BalanceRepoImpl) Create(
	ctx context.Context,
	model *models.Balance,
) (*models.Balance, error) {
	qu, _, err := goqu.
		Insert(balanceTName).
		Rows(model).
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

	var balance models.Balance
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&balance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert")
	}

	return &balance, nil
}

func (r *BalanceRepoImpl) Update(
	ctx context.Context,
	model *models.Balance,
) (*models.Balance, error) {
	qu, _, err := goqu.
		Update(balanceTName).
		Set(model).
		Where(goqu.C("balance_id").Eq(model.ID)).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	tx, err := r.repos.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to begin transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Error(context.Background(), "failed to rollback", err)
		}
	}()

	if _, err := tx.ExecContext(ctx, qu); err != nil {
		return nil, errors.Wrapf(err, "update balance error during execute query")
	}

	balance, err := r.Get(ctx, model.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "refresh balance error")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "failed to commit")
	}

	return balance, nil
}
