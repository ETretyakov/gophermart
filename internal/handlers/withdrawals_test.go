package handlers

import (
	"bytes"
	"context"
	"gophermart/internal/controllers"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/middlewares"
	"gophermart/internal/models"
	"gophermart/internal/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

type MockedWithdrawalsController struct{}

func (m *MockedWithdrawalsController) Create(
	ctx context.Context,
	schema *models.Withdrawal,
) (*models.Withdrawal, error) {
	switch {
	case schema.Order == "2312444431":
		return nil, exceptions.ErrWrongOrderNumber
	case schema.Order == "2357309406":
		return &models.Withdrawal{}, nil
	case schema.Order == "5754427317":
		return nil, exceptions.ErrBalanceIsNegative
	default:
		return &models.Withdrawal{}, nil
	}
}

func (m *MockedWithdrawalsController) UserWithdrawals(
	ctx context.Context,
	userID string,
) (*[]models.Withdrawal, error) {
	switch {
	case userID == "78b9e562-abda-4601-b3f5-ed7fbf089e8c":
		return &[]models.Withdrawal{{}}, nil
	case userID == "384325ad-f073-4ef3-9c7b-75ffbfef7399":
		return &[]models.Withdrawal{}, nil
	default:
		return nil, errors.New("user not found")
	}
}

func TestWithdrawalsHandlers_Create(t *testing.T) {
	controller := &MockedWithdrawalsController{}

	type fields struct {
		validator  validators.WithdrawalsValidator
		controller controllers.WithdrawalController
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
				validator:  validators.NewWithdrawalsValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestWithdrawalsHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/balance/withdraw",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"order": "2357309406", "sum": 2.0}`)
						}(),
					),
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Fail - Wrong order number",
			fields: fields{
				validator:  validators.NewWithdrawalsValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestWithdrawalsHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/balance/withdraw",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"order": "2312444431", "sum": 2.0}`)
						}(),
					),
				),
				wantStatusCode: 422,
			},
		},
		{
			name: "Test #3 Fail - Wrong body",
			fields: fields{
				validator:  validators.NewWithdrawalsValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestWithdrawalsHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/balance/withdraw",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"order": "2312444431", "sum": 2.0, "some": "some"}`)
						}(),
					),
				),
				wantStatusCode: 422,
			},
		},
		{
			name: "Test #4 Fail - Negative Balance",
			fields: fields{
				validator:  validators.NewWithdrawalsValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestWithdrawalsHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/balance/withdraw",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"order": "5754427317", "sum": 2.0}`)
						}(),
					),
				),
				wantStatusCode: 402,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &WithdrawalsHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}
			ctx := tt.args.r.Context()
			ctx = context.WithValue(ctx, middlewares.UserIDKey, tt.fields.userID)
			h.Create(tt.args.w, tt.args.r.WithContext(ctx))

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

func TestWithdrawalsHandlers_UserWithdrawals(t *testing.T) {
	controller := &MockedWithdrawalsController{}

	type fields struct {
		validator  validators.WithdrawalsValidator
		controller controllers.WithdrawalController
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
				validator:  validators.NewWithdrawalsValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestWithdrawalsHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/withdrawals",
					http.NoBody,
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Success - No content",
			fields: fields{
				validator:  validators.NewWithdrawalsValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestWithdrawalsHandlers"),
				userID:     "384325ad-f073-4ef3-9c7b-75ffbfef7399",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/withdrawals",
					http.NoBody,
				),
				wantStatusCode: 204,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &WithdrawalsHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}

			ctx := tt.args.r.Context()
			ctx = context.WithValue(ctx, middlewares.UserIDKey, tt.fields.userID)
			h.UserWithdrawals(tt.args.w, tt.args.r.WithContext(ctx))

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
