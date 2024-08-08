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
)

type MockedOrdersController struct{}

func (m *MockedOrdersController) Create(
	ctx context.Context,
	schema *models.Order,
) (*models.Order, error) {
	switch {
	case schema.Number == "2357309406":
		return &models.Order{}, exceptions.ErrOrderAlreadyAccepted
	case schema.Number == "5754427317":
		return &models.Order{}, exceptions.ErrOrderAlreadyRegistered
	default:
		return &models.Order{}, nil
	}
}

func (m *MockedOrdersController) UserOrders(
	ctx context.Context,
	userID string,
) (*[]models.Order, error) {
	switch {
	case userID == "78b9e562-abda-4601-b3f5-ed7fbf089e8c":
		return &[]models.Order{}, nil
	default:
		orders := make([]models.Order, 10)
		return &orders, nil
	}
}

func (m *MockedOrdersController) GetUserOrderByNumber(
	ctx context.Context,
	userID string,
	orderNumber string,
) (*models.Order, error) {
	switch {
	case orderNumber == "2357309406":
		return nil, exceptions.ErrOrderNotFound
	default:
		return &models.Order{}, nil
	}
}

func TestOrdersHandlers_Create(t *testing.T) {
	controller := &MockedOrdersController{}

	type fields struct {
		validator  validators.OrderValidator
		controller controllers.OrdersController
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
			name: "Test #1 Success - Accepted",
			fields: fields{
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/orders",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`4311889127`)
						}(),
					),
				),
				wantStatusCode: 202,
			},
		},
		{
			name: "Test #2 Success - Already Accepted",
			fields: fields{
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/orders",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`2357309406`)
						}(),
					),
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #3 Fail - Already Registered",
			fields: fields{
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/api/user/orders",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`5754427317`)
						}(),
					),
				),
				wantStatusCode: 409,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &OrdersHandlers{
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

func TestOrdersHandlers_UserOrders(t *testing.T) {
	controller := &MockedOrdersController{}

	type fields struct {
		validator  validators.OrderValidator
		controller controllers.OrdersController
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
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "2300c1c5-628f-4189-9a2d-47b8bb0cdfd0",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/orders",
					http.NoBody,
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Success - No content",
			fields: fields{
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/orders",
					http.NoBody,
				),
				wantStatusCode: 204,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &OrdersHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}
			ctx := tt.args.r.Context()
			ctx = context.WithValue(ctx, middlewares.UserIDKey, tt.fields.userID)
			h.UserOrders(tt.args.w, tt.args.r.WithContext(ctx))

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

func TestOrdersHandlers_UserOrderByNumber(t *testing.T) {
	controller := &MockedOrdersController{}

	type fields struct {
		validator  validators.OrderValidator
		controller controllers.OrdersController
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
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/orders/5754427317",
					http.NoBody,
				),
				wantStatusCode: 200,
			},
		},
		{
			name: "Test #2 Success - No content",
			fields: fields{
				validator:  validators.NewOrdersValidator(),
				controller: controller,
				logger:     log.NewHTTPLogger("TestOrdersHandlers"),
				userID:     "78b9e562-abda-4601-b3f5-ed7fbf089e8c",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/api/user/orders/2357309406",
					http.NoBody,
				),
				wantStatusCode: 204,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &OrdersHandlers{
				validator:  tt.fields.validator,
				controller: tt.fields.controller,
				logger:     tt.fields.logger,
			}
			ctx := tt.args.r.Context()
			ctx = context.WithValue(ctx, middlewares.UserIDKey, tt.fields.userID)
			h.UserOrderByNumber(tt.args.w, tt.args.r.WithContext(ctx))
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
