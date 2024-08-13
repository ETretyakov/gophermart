package handlers

import (
	"context"
	"gophermart/internal/controllers"
	"gophermart/internal/log"
	"gophermart/internal/middlewares"
	"gophermart/internal/models"
	"gophermart/internal/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

type MockedBalanceController struct{}

func (m *MockedBalanceController) GetForUser(
	ctx context.Context,
	userID string,
) (*models.Balance, error) {
	switch {
	case userID == "78b9e562-abda-4601-b3f5-ed7fbf089e8c":
		return &models.Balance{}, nil
	default:
		return nil, errors.New("user not found")
	}
}

func TestBalanceHandlers_GetForUser(t *testing.T) {
	controller := &MockedBalanceController{}

	type fields struct {
		validator  validators.BalanceValidator
		controller controllers.BalanceController
		logger     log.HTTPLogger
		userID     string
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
				validator:  validators.NewBalanceValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestBalanceHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/balance",
					http.NoBody,
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Unauthorized",
			fields: fields{
				validator:  validators.NewBalanceValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestBalanceHandlers"),
				userID:     "",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/balance",
					http.NoBody,
				),
				wantStatusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BalanceHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}

			ctx := tt.args.r.Context()
			ctx = context.WithValue(ctx, middlewares.UserIDKey, tt.fields.userID)
			h.GetForUser(tt.args.w, tt.args.r.WithContext(ctx))

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
