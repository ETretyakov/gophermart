package main

import (
	"context"
	"gophermart/internal/app"
	"gophermart/internal/config"
	"gophermart/internal/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error(ctx, "failed to read config", err)
	}

	log.InitDefault(cfg.LogLevel)

	// middlewares.SetSignKey(cfg.SignKey)

	if err := app.Run(ctx, cfg); err != nil {
		log.Error(ctx, "error running http server: %v", err)
	}
}
