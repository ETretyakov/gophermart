package controllers

import (
	"context"
	"gophermart/internal/models"
	"gophermart/internal/repository"

	"github.com/pkg/errors"
)

type BalanceControllerImpl struct {
	repos *repository.Repos
}

func NewBalanceController(repos *repository.Repos) *BalanceControllerImpl {
	return &BalanceControllerImpl{repos: repos}
}

func (c *BalanceControllerImpl) GetForUser(
	ctx context.Context,
	userID string,
) (*models.Balance, error) {
	balance, err := c.repos.BalanceRepo.GetOrCreateForUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get balance")
	}

	return balance, nil
}
