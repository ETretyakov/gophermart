package validators

import (
	"context"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type OrdersValidatorImpl struct {
	validate *validator.Validate
}

func NewOrdersValidator() *OrdersValidatorImpl {
	return &OrdersValidatorImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *OrdersValidatorImpl) ValidateOrderCreate(
	userID string,
	body io.ReadCloser,
) (*models.Order, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read body")
	}
	defer func() {
		err := body.Close()
		if err != nil {
			log.Error(context.Background(), "failed to read body", err)
		}
	}()

	rawOrderNumber := string(bodyBytes)
	if err := goluhn.Validate(rawOrderNumber); err != nil {
		return nil, exceptions.ErrWrongOrderNumber
	}

	order := models.NewOrder(userID, rawOrderNumber)

	return order, nil
}

func (v *OrdersValidatorImpl) ValidateOrderFromPath(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	number, ok := vars["number"]
	if !ok {
		return "", errors.New("failed to retrieve number")
	}

	if err := goluhn.Validate(number); err != nil {
		return "", exceptions.ErrWrongOrderNumber
	}

	return number, nil
}
