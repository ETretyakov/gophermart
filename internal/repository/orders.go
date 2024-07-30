package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"gophermart/internal/types"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

type OrdersRepoImpl struct {
	repos *Repos
}

func NewOrdersRepoImpl(repos *Repos) *OrdersRepoImpl {
	return &OrdersRepoImpl{repos: repos}
}

func (r *OrdersRepoImpl) Create(
	ctx context.Context,
	model *models.Order,
) (*models.Order, error) {
	qu, _, err := goqu.
		Insert(ordersTName).
		Rows(model).
		Returning(
			"order_id",
			"user_id",
			"number",
			"status",
			"accrual",
			"uploaded_at",
		).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var order models.Order
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&order)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert")
	}

	return &order, nil
}

func (r *OrdersRepoImpl) GetByNumber(
	ctx context.Context,
	orderNumber uint64,
) (*models.Order, error) {
	qu, _, err := goqu.
		Select(&models.Order{}).
		From(ordersTName).
		Where(goqu.C("number").Eq(orderNumber)).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var order models.Order
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrOrderNotFound
		} else {
			return nil, errors.Wrapf(err, "failed to execute query")
		}
	}

	return &order, nil
}

func (r *OrdersRepoImpl) GetUserOrderByNumber(
	ctx context.Context,
	userID string,
	orderNumber uint64,
) (*models.Order, error) {
	qu, _, err := goqu.
		Select(&models.Order{}).
		From(ordersTName).
		Where(
			goqu.C("number").Eq(orderNumber),
			goqu.C("user_id").Eq(userID),
		).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	var order models.Order
	err = r.repos.DB.QueryRowxContext(ctx, qu).StructScan(&order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrOrderNotFound
		} else {
			return nil, errors.Wrapf(err, "failed to execute query")
		}
	}

	return &order, nil
}

func (r *OrdersRepoImpl) changeStatus(
	ctx context.Context,
	orderIDs []string,
	status types.OrderStatus,
) (bool, error) {
	qu, _, err := goqu.
		Update(ordersTName).
		Set(map[string]interface{}{"status": status.String()}).
		Where(
			goqu.C("order_id").In(orderIDs),
		).
		ToSQL()
	if err != nil {
		return false, errors.Wrapf(err, "failed to build query")
	}

	_, err = r.repos.DB.ExecContext(ctx, qu)
	if err != nil {
		return false, errors.Wrapf(err, "failed to update")
	}

	return true, nil
}

func (r *OrdersRepoImpl) MarkAsProcessing(
	ctx context.Context,
	orderIDs []string,
) (bool, error) {
	status, err := r.changeStatus(ctx, orderIDs, types.OrderProcessing)
	if err != nil {
		return false, errors.Wrapf(err, "failed to change status")
	}

	return status, nil
}

func (r *OrdersRepoImpl) MarkAsInvalid(
	ctx context.Context,
	orderIDs []string,
) (bool, error) {
	status, err := r.changeStatus(ctx, orderIDs, types.OrderInvalid)
	if err != nil {
		return false, errors.Wrapf(err, "failed to change status")
	}

	return status, nil
}

func (r *OrdersRepoImpl) Accrue(
	ctx context.Context,
	record models.AccrueRecord,
) (bool, error) {
	_, err := r.repos.BalanceRepo.GetOrCreateForUser(ctx, record.UserID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get user balance")
	}

	// Setup transaction
	tx, err := r.repos.DB.BeginTxx(ctx, nil)
	if err != nil {
		return false, errors.Wrapf(err, "failed to begin transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Error(context.Background(), "failed to rollback", err)
		}
	}()

	// Update order
	qu, _, err := goqu.
		Update(ordersTName).
		Set(
			map[string]interface{}{
				"accrual": record.Amount,
				"status":  types.OrderProcessed.String(),
			},
		).
		Where(goqu.C("number").Eq(record.Number)).
		Returning(
			"order_id",
			"user_id",
			"number",
			"status",
			"accrual",
			"uploaded_at",
		).
		ToSQL()
	if err != nil {
		return false, errors.Wrapf(err, "failed to build query")
	}

	var order models.Order
	err = tx.QueryRowxContext(ctx, qu).StructScan(&order)
	if err != nil {
		return false, errors.Wrapf(err, "failed to update order")
	}

	// Making balance update
	qu, _, err = goqu.
		Update(balanceTName).
		Set(
			map[string]interface{}{
				"current": goqu.L(fmt.Sprintf("current + %f", record.Amount)),
			},
		).
		Where(goqu.C("user_id").Eq(record.UserID)).
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
		return false, errors.Wrapf(err, "failed to build query")
	}
	if _, err := tx.ExecContext(ctx, qu); err != nil {
		return false, errors.Wrapf(err, "update balance error during execute query")
	}

	if err := tx.Commit(); err != nil {
		return false, errors.Wrapf(err, "failed to commti on accrue")
	}

	return true, nil
}

func (r *OrdersRepoImpl) UserOrders(
	ctx context.Context,
	userID string,
) (*[]models.Order, error) {
	qu, _, err := goqu.
		Select(&models.Order{}).
		From(ordersTName).
		Where(goqu.C("user_id").Eq(userID)).
		Order(goqu.I("uploaded_at").Asc()).
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query")
	}

	rows, err := r.repos.DB.QueryxContext(ctx, qu)
	if err != nil {
		return nil, errors.Wrapf(err, "read orders error during querying")
	}

	orders := []models.Order{}
	for rows.Next() {
		order := models.Order{}
		err := rows.StructScan(&order)
		if err != nil {
			return nil, errors.Wrapf(err, "read orders error during scan rows")
		}
		orders = append(orders, order)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("read orders error during querying: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(ctx, "failed to close rows", err)
		}
	}()

	return &orders, nil
}
