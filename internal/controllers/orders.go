package controllers

import (
	"context"
	"gophermart/internal/exceptions"
	"gophermart/internal/models"
	"gophermart/internal/pipelines"
	"gophermart/internal/repository"

	"github.com/pkg/errors"
)

type AccrualPipeline interface {
	RegisterOrder(order *models.Order)
}

type OrdersControllerImpl struct {
	repos *repository.Repos
}

func NewOrdersController(
	repos *repository.Repos,
) *OrdersControllerImpl {
	return &OrdersControllerImpl{repos: repos}
}

func (c *OrdersControllerImpl) Create(
	ctx context.Context,
	schema *models.Order,
) (*models.Order, error) {
	existingOrder, err := c.repos.OrdersRepo.GetByNumber(ctx, schema.Number)
	if err != nil {
		if !errors.Is(err, exceptions.ErrOrderNotFound) {
			return nil, errors.Wrapf(err, "failed to get order")
		}
	}
	if existingOrder != nil {
		if existingOrder.UserID == schema.UserID {
			return nil, exceptions.ErrOrderAlreadyAccepted
		} else {
			return nil, exceptions.ErrOrderAlreadyRegistered
		}
	}

	order, err := c.repos.OrdersRepo.Create(ctx, schema)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create order")
	}

	pipelines.AccrualPipeline.RegisterOrder(order)

	return order, nil
}

func (c *OrdersControllerImpl) UserOrders(
	ctx context.Context,
	userID string,
) (*[]models.Order, error) {
	orders, err := c.repos.OrdersRepo.UserOrders(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get orders")
	}

	return orders, nil
}

func (c *OrdersControllerImpl) GetUserOrderByNumber(
	ctx context.Context,
	userID string,
	orderNumber uint64,
) (*models.Order, error) {
	order, err := c.repos.OrdersRepo.GetUserOrderByNumber(
		ctx,
		userID,
		orderNumber,
	)
	if err != nil {
		if errors.Is(err, exceptions.ErrOrderNotFound) {
			return nil, exceptions.ErrOrderNotFound
		} else {
			return nil, errors.Wrapf(err, "failed to get order")
		}
	}

	return order, nil
}
