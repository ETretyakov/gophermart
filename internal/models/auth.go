package models

import (
	"gophermart/internal/crypto"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AuthUser struct {
	ID        string  `db:"user_id"`
	Login     string  `db:"login"`
	Password  string  `db:"hashed_password"`
	CreatedAt string  `db:"created_at"`
	UpdatedAt string  `db:"updated_at"`
	DeletedAt *string `db:"deleted_at"`
}

func NewAuthUser(
	login string,
	password string,
) (*AuthUser, error) {
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hash password")
	}

	return &AuthUser{
		ID:        uuid.NewString(),
		Login:     login,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

type UserRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
