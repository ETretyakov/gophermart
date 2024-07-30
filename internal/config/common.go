package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Environment string

const (
	EnvProd  Environment = "prod"
	EnvDev   Environment = "dev"
	EnvLocal Environment = "local"
)

type Postgres struct {
	DSN             string        `env:"URI"              envDefault:""`
	MaxOpenConn     int           `env:"MAX_OPEN_CONN"    envDefault:"10"`
	IdleConn        int           `env:"MAX_IDLE_CONN"    envDefault:"10"`
	PingInterval    time.Duration `env:"DURATION"         envDefault:"5s"`
	MigrationFolder string        `env:"MIGRATION_FOLDER" envDefault:"./migrations"`
}

type Security struct {
	JWTSecretKey string        `env:"JWT_SECRET_KEY" envDefault:"unsecured_key"`
	JWTExpire    time.Duration `env:"JWT_EXPIRE" envDefault:"12h"`
}

type Config struct {
	AppEnv                         Environment   `env:"APP_ENVIRONMENT" envDefault:"local" flag:"mode" flagShort:"m" flagDescription:"environment"`
	HTTPAddress                    string        `env:"APP_ADDRESS" envDefault:"localhost:8081" flag:"address" flagShort:"a" flagDescription:"http address"`
	LogLevel                       string        `env:"LOG_LEVEL" envDefault:"info" flag:"log_level" flagShort:"l" flagDescription:"level for logging"`
	LogFile                        string        `env:"LOG_FILE" envDefault:"logs/logs.jsonl" flag:"log_file"  flagShort:"w" flagDescription:"filepath for logs"`
	Postgres                       Postgres      `envPrefix:"DATABASE_" flag:"pg_dsn" flagShort:"d" flagDescription:"database dsn"`
	Security                       Security      `envPrefix:"SECURITY_" flag:"jwt_secret_key" flagShort:"j" flagDescription:"jwt secret key"`
	AccrualBaseURL                 string        `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"localhost:8080" flag:"accrual_address" flagShort:"r" flagDescription:"accrual address"`
	AccrualRetryCount              int           `env:"ACCRUAL_RETRY_COUNT" envDefault:"3"`
	AccrualRetryWaitTime           time.Duration `env:"ACCRUAL_RETRY_WAIT_TIME" envDefault:"1s"`
	AccrualRetryMaxWaitTime        time.Duration `env:"ACCRUAL_RETRY_MAX_WAIT_TIME" envDefault:"10s"`
	AccrualPipelineBufferSize      int           `env:"ACCRUAL_PIPELINE_BUFFER_SIZE" envDefault:"10"`
	AccrualPipelineNumberOfWorkers int           `env:"ACCRUAL_PIPELINE_NUMBER_OF_WORKERS" envDefault:"10"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse server envs: %w", err)
	}

	ParseFlags(cfg)

	return cfg, nil
}
