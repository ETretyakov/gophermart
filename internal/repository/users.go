package repository

import (
	"context"
	"gophermart/internal/log"
	"gophermart/internal/models"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

type UsersRepoImpl struct {
	repos *Repos
}

func NewUsersRepo(repos *Repos) *UsersRepoImpl {
	return &UsersRepoImpl{repos: repos}
}

func (r *UsersRepoImpl) Create(
	ctx context.Context,
	model *models.AuthUser,
) (*models.User, error) {
	// Setup transaction
	tx, err := r.repos.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to begin transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Error(context.Background(), "failed to rollback", err)
		}
	}()

	// Create User
	qu, _, err := goqu.
		Insert(usersTName).
		Rows(model).
		Returning(
			"user_id",
			"login",
			"created_at",
			"updated_at",
			"deleted_at",
		).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var user models.User
	err = tx.QueryRowxContext(ctx, qu).StructScan(&user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert")
	}

	// Create Balance
	qu, _, err = goqu.
		Insert(balanceTName).
		Rows(models.NewBalance(user.ID)).
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
	_, err = tx.ExecContext(ctx, qu)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert balance")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "failed to commti on user create")
	}

	return &user, nil
}
