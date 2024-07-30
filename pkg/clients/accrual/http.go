package accrual

import (
	"context"
	"fmt"
	"gophermart/internal/log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type AccrualHTTPClient struct {
	client  *resty.Client
	baseURL string
}

func NewAccrualClient(
	ctx context.Context,
	baseURL string,
	retryCount int,
	retryWaitTime time.Duration,
	retryMaxWaitTime time.Duration,
) *AccrualHTTPClient {
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}

	c := &AccrualHTTPClient{
		client:  resty.New(),
		baseURL: baseURL,
	}

	c.client.
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetRetryMaxWaitTime(retryMaxWaitTime)

	return c
}

func (c *AccrualHTTPClient) CreateGoods(
	ctx context.Context,
	schema GoodsCreate,
) error {
	request := c.client.R().
		SetContext(ctx).
		SetBody(&schema)

	resp, err := request.Post(c.baseURL + "/api/goods")
	if err != nil {
		return errors.Wrapf(err, "failed to create goods")
	}

	if resp.IsError() {
		return errors.Wrapf(
			err,
			"failed to create goods: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		)
	}

	log.Info(
		ctx,
		fmt.Sprintf(
			"reponse: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		),
	)

	return nil
}

func (c *AccrualHTTPClient) CreateOrder(
	ctx context.Context,
	schema OrderCreate,
) error {
	request := c.client.R().
		SetContext(ctx).
		SetBody(&schema)

	resp, err := request.Post(c.baseURL + "/api/orders")
	if err != nil {
		return errors.Wrapf(err, "failed to create order")
	}

	if resp.IsError() {
		return errors.Wrapf(
			err,
			"failed to create order: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		)
	}

	log.Info(
		ctx,
		fmt.Sprintf(
			"reponse: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		),
	)

	return nil
}

func (c *AccrualHTTPClient) GetOrder(
	ctx context.Context,
	order uint64,
) (*OrderRead, error) {
	var orderModel OrderRead

	request := c.client.R().
		SetContext(ctx).
		SetResult(&orderModel)

	resp, err := request.Get(fmt.Sprintf("%s/api/orders/%d", c.baseURL, order))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get order")
	}

	if resp.IsError() {
		return nil, errors.Wrapf(
			err,
			"failed to get order: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		)
	}

	log.Info(
		ctx,
		fmt.Sprintf(
			"reponse: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		),
	)

	return &orderModel, nil
}
