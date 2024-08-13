package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gophermart/internal/bootstrap"
	"gophermart/internal/closer"
	"gophermart/internal/config"
	"gophermart/internal/crypto"
	"gophermart/internal/handlers"
	"gophermart/internal/http"
	"gophermart/internal/log"
	"gophermart/internal/pipelines"
	"gophermart/internal/repository"

	"github.com/jmoiron/sqlx"
)

func Run(ctx context.Context, cfg *config.Config) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	crypto.InitJWTSigner(cfg.Security.JWTSecretKey, cfg.Security.JWTExpire)

	// Database setup
	var db *sqlx.DB
	if cfg.Postgres.DSN != "" {
		db, err = bootstrap.InitDB(ctx, &cfg.Postgres)
		if err != nil {
			log.Fatal(ctx, "failed to init db", err)
		}
	}

	repos, err := repository.NewRepos(
		ctx,
		db,
	)
	if err != nil {
		log.Fatal(ctx, "failed to init repos", err)
	}

	// Pipelines
	pipelines.InitAccrualPipeline(
		ctx,
		cfg.AccrualBaseURL,
		cfg.AccrualRetryCount,
		cfg.AccrualRetryWaitTime,
		cfg.AccrualRetryMaxWaitTime,
		repos,
		cfg.AccrualPipelineBufferSize,
		cfg.AccrualPipelineNumberOfWorkers,
	)

	pipelines.AccrualPipeline.Start(ctx)

	// Handlers bindings
	healthHandlers := handlers.NewHealthHandlers(repos)
	authHandlers := handlers.NewAuthHandlers(repos)
	balanceHandlers := handlers.NewBalanceHandlers(repos)
	ordersHandlers := handlers.NewOrdersHandlers(repos)
	withdrawalsHandlers := handlers.NewWithdrawalsHandlers(repos)

	httpServer := http.New(
		cfg,
		healthHandlers,
		authHandlers,
		balanceHandlers,
		ordersHandlers,
		withdrawalsHandlers,
	)

	// Server start
	httpServer.Start(ctx)

	healthHandlers.SetLiveness(true)
	healthHandlers.SetReadiness(true)

	gracefulShutDown(ctx, cancel)

	return nil
}

func gracefulShutDown(ctx context.Context, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	errorMessage := fmt.Sprintf("%s %v - %s", "Received shutdown signal:", sig, "Graceful shutdown done")
	log.Info(ctx, errorMessage)
	cancel()
	closer.CloseAll()
}
