package validators

import "github.com/go-playground/validator/v10"

type BalanceValidatorImpl struct {
	validate *validator.Validate
}

func NewBalanceValidator() *BalanceValidatorImpl {
	return &BalanceValidatorImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}
