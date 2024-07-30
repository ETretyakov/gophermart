package validators

import (
	"encoding/json"
	"gophermart/internal/exceptions"
	"gophermart/internal/models"
	"io"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type WithdrawalsValidatorImpl struct {
	validate *validator.Validate
}

func NewWithdrawalsValidator() *WithdrawalsValidatorImpl {
	return &WithdrawalsValidatorImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *WithdrawalsValidatorImpl) ValidateOrderCreate(
	userID string,
	body io.ReadCloser,
) (*models.Withdrawal, error) {
	withdrawalCreate := &models.WithdrawalCreate{}

	err := json.NewDecoder(body).Decode(withdrawalCreate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse user register json")
	}

	if err := goluhn.Validate(withdrawalCreate.Order); err != nil {
		return nil, exceptions.ErrWrongOrderNumber
	}

	withdrawal := models.NewWithdrawal(
		userID,
		withdrawalCreate.Order,
		withdrawalCreate.Sum,
	)

	return withdrawal, nil
}
