package pipelines

import (
	"context"
	"gophermart/internal/repository"
	"gophermart/pkg/clients/accrual"
	"time"
)

var AccrualPipeline AccrualPipelineImpl

func InitAccrualPipeline(
	ctx context.Context,
	accrualBaseURL string,
	accrualRetryCount int,
	accrualRetryWaitTime time.Duration,
	accrualRetryMaxWaitTime time.Duration,
	repos *repository.Repos,
	bufferSize int,
	numberOfWorkers int,
) {
	AccrualPipeline = *NewAccrualPipeline(
		repository.NewOrdersRepoImpl(repos),
		accrual.NewAccrualClient(
			ctx,
			accrualBaseURL,
			accrualRetryCount,
			accrualRetryWaitTime,
			accrualRetryMaxWaitTime,
		),
		bufferSize,
		numberOfWorkers,
	)
}
