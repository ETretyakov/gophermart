package pipelines

import (
	"context"
	"fmt"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"gophermart/internal/types"
	"gophermart/pkg/clients/accrual"
)

type OrdersRepo interface {
	MarkAsProcessing(ctx context.Context, orderIDs []string) (bool, error)
	MarkAsInvalid(ctx context.Context, orderIDs []string) (bool, error)
	Accrue(ctx context.Context, record models.AccrueRecord) (bool, error)
}

type AccrualClient interface {
	GetOrder(ctx context.Context, order string) (*accrual.OrderRead, error)
}

type AccrualPipelineImpl struct {
	client          AccrualClient
	ordersRepo      OrdersRepo
	preprocessingCh chan models.Order
	processingCh    chan models.AccrueRecord
	numberOfWorkers int
}

func NewAccrualPipeline(
	ordersRepo OrdersRepo,
	client AccrualClient,
	bufferSize int,
	numberOfWorkers int,
) *AccrualPipelineImpl {
	return &AccrualPipelineImpl{
		client:          client,
		ordersRepo:      ordersRepo,
		preprocessingCh: make(chan models.Order, bufferSize),
		processingCh:    make(chan models.AccrueRecord, bufferSize),
		numberOfWorkers: numberOfWorkers,
	}
}

func (p *AccrualPipelineImpl) RegisterOrder(order *models.Order) {
	p.preprocessingCh <- *order
}

func (p *AccrualPipelineImpl) preprocessingWorker(ctx context.Context, workerID int) {
	log.Info(ctx, fmt.Sprintf("starting preprocessing Worker №%d", workerID))

	for {
		select {
		case order := <-p.preprocessingCh:
			log.Info(
				ctx,
				fmt.Sprintf(
					"preprocessing Worker №%d: got order=%s",
					workerID,
					order.Number,
				),
			)

			orderRead, err := p.client.GetOrder(ctx, order.Number)
			if err != nil {
				log.Error(ctx, "failed to get order info", err)
			}

			log.Info(
				ctx,
				fmt.Sprintf(
					"preprocessing Worker №%d: retrieved status for order=%s - %+v",
					workerID,
					order.Number,
					orderRead,
				),
			)

			switch {
			case orderRead.Status == string(types.OrderProcessed):
				accrueRecord := models.AccrueRecord{
					UserID: order.UserID,
					Number: order.Number,
					Amount: orderRead.Accrual,
				}

				p.processingCh <- accrueRecord

				log.Info(
					ctx,
					fmt.Sprintf(
						"preprocessing Worker №%d: order=%s sent to processing",
						workerID,
						order.Number,
					),
				)
			case orderRead.Status == string(types.OrderProcessing):
				_, err := p.ordersRepo.MarkAsProcessing(ctx, []string{order.ID})
				if err != nil {
					log.Error(ctx, "failed to mark processing order", err)
				}

				log.Info(
					ctx,
					fmt.Sprintf(
						"preprocessing Worker №%d: marked order=%s as processing",
						workerID,
						order.Number,
					),
				)
			case orderRead.Status == string(types.OrderInvalid):
				_, err := p.ordersRepo.MarkAsInvalid(ctx, []string{order.ID})
				if err != nil {
					log.Error(ctx, "failed to mark invalid order", err)
				}

				log.Info(
					ctx,
					fmt.Sprintf(
						"preprocessing Worker №%d: order=%s marked as invalid",
						workerID,
						order.Number,
					),
				)
			default:
				log.Info(
					ctx,
					fmt.Sprintf(
						"preprocessing Worker №%d: order=%s status unchanged",
						workerID,
						order.Number,
					),
				)
			}
		case <-ctx.Done():
			log.Info(ctx, fmt.Sprintf("preprocessing Worker №%d shutdown", workerID))
			return
		}
	}
}

func (p *AccrualPipelineImpl) processingWorker(ctx context.Context, workerID int) {
	log.Info(ctx, fmt.Sprintf("starting processing Worker №%d", workerID))

	for {
		select {
		case accrueRecord := <-p.processingCh:
			log.Info(
				ctx,
				fmt.Sprintf(
					"processing Worker №%d: got order=%s",
					workerID,
					accrueRecord.Number,
				),
			)

			_, err := p.ordersRepo.Accrue(ctx, accrueRecord)
			if err != nil {
				log.Error(ctx, "failed to accrue order", err)
			}

			log.Info(
				ctx,
				fmt.Sprintf(
					"processing Worker №%d: order=%s processed",
					workerID,
					accrueRecord.Number,
				),
			)
		case <-ctx.Done():
			log.Info(ctx, fmt.Sprintf("processing Worker №%d shutdown", workerID))
			return
		}
	}
}

func (p *AccrualPipelineImpl) Start(ctx context.Context) {
	log.Info(ctx, fmt.Sprintf("Starting %d workers", p.numberOfWorkers))
	for i := 1; i <= p.numberOfWorkers; i++ {
		go p.preprocessingWorker(ctx, i)
		go p.processingWorker(ctx, i)
	}
}
