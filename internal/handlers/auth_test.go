package handlers

import (
	"bytes"
	"context"
	"gophermart/internal/controllers"
	"gophermart/internal/crypto"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"gophermart/internal/validators"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MockedAuthController struct{}

func (m *MockedAuthController) Register(
	ctx context.Context,
	schema *models.AuthUser,
) (*models.User, error) {
	if schema.Login == "test3" {
		return nil, exceptions.ErrLoginAlreadyTaken
	}

	authUser, err := models.NewAuthUser(schema.Login, schema.Password)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate auth user")
	}
	return &models.User{
		ID:        authUser.ID,
		Login:     authUser.Login,
		CreatedAt: authUser.CreatedAt,
		UpdatedAt: authUser.UpdatedAt,
		DeletedAt: authUser.DeletedAt,
	}, nil
}

func (m *MockedAuthController) Login(
	ctx context.Context,
	schema *models.UserLogin,
) (string, error) {
	crypto.InitJWTSigner("test_secure_key", 10*time.Second)
	if schema.Login == "test2" {
		return "", exceptions.ErrNotAuthorised
	}
	if schema.Login == "test3" {
		return "", exceptions.ErrUserNotFound
	}

	token, err := crypto.JWT.GetToken(uuid.NewString())
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate token")
	}

	return token, nil
}

func TestAuthHandlers_Register(t *testing.T) {
	controller := &MockedAuthController{}

	type fields struct {
		validator  validators.AuthValidator
		controller controllers.AuthController
		logger     log.HTTPLogger
	}
	type args struct {
		w              http.ResponseWriter
		r              *http.Request
		wantStatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test #1 Success",
			fields: fields{
				validator:  validators.NewAuthValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestAuthHandlers"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/register",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"login": "test1", "password": "12345Hello!"}`)
						}(),
					),
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Wrong Password",
			fields: fields{
				validator:  validators.NewAuthValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestAuthHandlers"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/register",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"login": "test2", "password": "123321"}`)
						}(),
					),
				),
				wantStatusCode: 400,
			},
		},
		{
			name: "Test #3 Login already taken",
			fields: fields{
				validator:  validators.NewAuthValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestAuthHandlers"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/register",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"login": "test3", "password": "12345Hello!"}`)
						}(),
					),
				),
				wantStatusCode: 409,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}

			h.Register(tt.args.w, tt.args.r)

			if r, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if r.Code != tt.args.wantStatusCode {
					t.Errorf(
						"status codes are different: got=%d want=%d",
						r.Code,
						tt.args.wantStatusCode,
					)
				}
			} else {
				t.Errorf("got different from *httptest.ResponseRecorder struct: %+v", r)
			}
		})
	}
}

func TestAuthHandlers_Login(t *testing.T) {
	controller := &MockedAuthController{}

	type fields struct {
		validator  validators.AuthValidator
		controller controllers.AuthController
		logger     log.HTTPLogger
	}
	type args struct {
		w              http.ResponseWriter
		r              *http.Request
		wantStatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test #1 Success",
			fields: fields{
				validator:  validators.NewAuthValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestAuthHandlers"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/login",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"login": "test1", "password": "12345Hello!"}`)
						}(),
					),
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Not authorised",
			fields: fields{
				validator:  validators.NewAuthValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestAuthHandlers"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/login",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"login": "test2", "password": "12345Hello!"}`)
						}(),
					),
				),
				wantStatusCode: 401,
			},
		},
		{
			name: "Test #3 User not found",
			fields: fields{
				validator:  validators.NewAuthValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestAuthHandlers"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/login",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"login": "test3", "password": "12345Hello!"}`)
						}(),
					),
				),
				wantStatusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}
			h.Login(tt.args.w, tt.args.r)

			if r, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if r.Code != tt.args.wantStatusCode {
					t.Errorf(
						"status codes are different: got=%d want=%d",
						r.Code,
						tt.args.wantStatusCode,
					)
				}
			} else {
				t.Errorf("got different from *httptest.ResponseRecorder struct: %+v", r)
			}
		})
	}
}
