package repository

import (
	"context"
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
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert")
	}

	return &user, nil
}
