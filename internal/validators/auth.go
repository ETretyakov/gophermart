package validators

import (
	"encoding/json"
	"gophermart/internal/models"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type AuthValidatorImpl struct {
	validate *validator.Validate
}

func NewAuthValidator() *AuthValidatorImpl {
	return &AuthValidatorImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *AuthValidatorImpl) ValidateUserRegister(body io.ReadCloser) (*models.AuthUser, error) {
	userRegister := &models.UserRegister{}

	err := json.NewDecoder(body).Decode(userRegister)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse user register json")
	}

	// Check if neccessary symbols in place
	if verified := VerifyPassword(userRegister.Password); !verified {
		return nil, errors.New("password don't match minimal requirements")
	}

	user, err := models.NewAuthUser(
		userRegister.Login,
		userRegister.Password,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get new auth user")
	}

	return user, nil
}

func (v *AuthValidatorImpl) ValidateUserLogin(body io.ReadCloser) (*models.UserLogin, error) {
	userLogin := &models.UserLogin{}

	err := json.NewDecoder(body).Decode(userLogin)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse user login json")
	}

	return userLogin, nil
}
