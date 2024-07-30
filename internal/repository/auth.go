package repository

import (
	"context"
	"database/sql"
	"gophermart/internal/exceptions"
	"gophermart/internal/models"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

type AuthRepoImpl struct {
	repos *Repos
}

func NewAuthRepo(repos *Repos) *AuthRepoImpl {
	return &AuthRepoImpl{repos: repos}
}

func (r *AuthRepoImpl) UserExists(
	ctx context.Context,
	login string,
) (bool, error) {
	qu, _, err := goqu.
		Select(goqu.C("user_id")).
		From(usersTName).
		Where(goqu.C("login").Eq(login)).
		ToSQL()
	if err != nil {
		return false, errors.Wrapf(err, "failed to build query")
	}

	var userID string
	err = r.repos.DB.QueryRowxContext(ctx, qu).Scan(&userID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, nil
}

func (r *AuthRepoImpl) UserAuth(
	ctx context.Context,
	login string,
) (*models.AuthUser, error) {
	qu, _, err := goqu.
		Select(
			goqu.C("user_id"),
			goqu.C("login"),
			goqu.C("hashed_password"),
			goqu.C("created_at"),
			goqu.C("updated_at"),
			goqu.C("deleted_at"),
		).
		From(usersTName).
		Where(goqu.C("login").Eq(login)).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var authUser models.AuthUser
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&authUser)
	switch {
	case err == nil:
		return &authUser, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, exceptions.ErrUserNotFound
	default:
		return nil, errors.Wrapf(err, "failed to query database")
	}
}
