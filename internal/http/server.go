package http

import (
	"context"
	"gophermart/internal/closer"
	"gophermart/internal/config"
	"gophermart/internal/handlers"
	"gophermart/internal/log"
	"net/http"

	"github.com/pkg/errors"
)

type Server struct {
	cfg         *config.Config
	srv         *http.Server
	health      *handlers.HealthHandlers
	auth        *handlers.AuthHandlers
	balance     *handlers.BalanceHandlers
	orders      *handlers.OrdersHandlers
	withdrawals *handlers.WithdrawalsHandlers
}

func New(
	cfg *config.Config,
	healthHandlers *handlers.HealthHandlers,
	authHandlers *handlers.AuthHandlers,
	balanceHandlers *handlers.BalanceHandlers,
	ordersHandlers *handlers.OrdersHandlers,
	withdrawalsHandlers *handlers.WithdrawalsHandlers,
) *Server {
	srv := &http.Server{
		Addr: cfg.HTTPAddress,
	}

	return &Server{
		cfg:         cfg,
		srv:         srv,
		health:      healthHandlers,
		auth:        authHandlers,
		balance:     balanceHandlers,
		orders:      ordersHandlers,
		withdrawals: withdrawalsHandlers,
	}
}

func (s *Server) Start(ctx context.Context) {
	s.srv.Handler = s.setupRoutes()

	go func() {
		log.Info(ctx, "starting listening http srv at "+s.cfg.HTTPAddress)
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(ctx, "error start http srv, err: %+v", err)
		}
	}()

	closer.Add(s.Close)
}

func (s *Server) Close() error {
	ctx := context.TODO()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error(ctx, "error stop http srv, err", err)
		return errors.Wrapf(err, "failed to shutdown server")
	}

	log.Info(ctx, "http server shutdown done")

	return nil
}
