package controllers

import (
	"context"
	"gophermart/internal/exceptions"
	"gophermart/internal/models"
	"gophermart/internal/repository"

	"github.com/pkg/errors"
)

type WithdrawalsControllerImpl struct {
	repos *repository.Repos
}

func NewWithdrawalsController(repos *repository.Repos) *WithdrawalsControllerImpl {
	return &WithdrawalsControllerImpl{repos: repos}
}

func (c *WithdrawalsControllerImpl) Create(
	ctx context.Context,
	schema *models.Withdrawal,
) (*models.Withdrawal, error) {
	withdrawal, err := c.repos.WithdrawalsRepo.Create(ctx, schema)
	if err != nil {
		if errors.Is(err, exceptions.ErrBalanceIsNegative) {
			return nil, exceptions.ErrBalanceIsNegative
		} else {
			return nil, errors.Wrapf(err, "failed to create withdrawal record")
		}
	}

	return withdrawal, nil
}

func (c *WithdrawalsControllerImpl) UserWithdrawals(
	ctx context.Context,
	userID string,
) (*[]models.Withdrawal, error) {
	withdrawals, err := c.repos.WithdrawalsRepo.UserWithdrawals(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get orders")
	}

	return withdrawals, nil
}
