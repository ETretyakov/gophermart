package controllers

import (
	"context"
	"fmt"
	"gophermart/internal/crypto"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"gophermart/internal/repository"

	"github.com/pkg/errors"
)

type AuthControllerImpl struct {
	repos *repository.Repos
}

func NewAuthController(repos *repository.Repos) *AuthControllerImpl {
	return &AuthControllerImpl{repos: repos}
}

func (c *AuthControllerImpl) Register(
	ctx context.Context,
	schema *models.AuthUser,
) (*models.User, error) {
	exists, err := c.repos.AuthRepo.UserExists(
		ctx,
		schema.Login,
	)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to check if login exists: %s",
			schema.Login,
		)
	}

	if exists {
		return nil, exceptions.ErrLoginAlreadyTaken
	}

	user, err := c.repos.UsersRepo.Create(ctx, schema)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed create user record",
		)
	}

	return user, nil
}

func (c *AuthControllerImpl) Login(
	ctx context.Context,
	schema *models.UserLogin,
) (string, error) {
	authUser, err := c.repos.AuthRepo.UserAuth(ctx, schema.Login)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrUserNotFound):
			return "", exceptions.ErrUserNotFound
		default:
			return "", errors.Wrapf(
				err,
				"failed retrieve hashed password",
			)
		}
	}

	if verified := crypto.CheckPasswordHash(
		schema.Password,
		authUser.Password,
	); !verified {
		return "", exceptions.ErrNotAuthorised
	}

	token, err := crypto.JWT.GetToken(authUser.ID)
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("failed to get token: %s", err))
		return "", exceptions.ErrNotAuthorised
	}

	return token, nil
}
