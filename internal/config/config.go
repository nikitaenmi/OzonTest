package config

import (
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	DB       DatabaseConfig
	Server   ServerConfig
	CBR      CBRConfig
	Payment  PaymentConfig
	LoadTest LoadTestConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	GRPCPort    string        `env:"GRPC_PORT"`
	HTTPPort    string        `env:"HTTP_PORT"`
	GRPCTimeout time.Duration `env:"GRPC_TIMEOUT"`
}

type CBRConfig struct {
	Timeout int `env:"CBR_TIMEOUT"`
}

type PaymentConfig struct {
	MaxAmountRUB float64 `env:"PAYMENT_MAX_AMOUNT_RUB"`
}

type LoadTestConfig struct {
	TotalRequests int           `env:"LOAD_TEST_REQUESTS"`
	Concurrency   int           `env:"LOAD_TEST_CONCURRENCY"`
	Timeout       time.Duration `env:"LOAD_TEST_TIMEOUT"`
	BaseURL       string        `env:"LOAD_TEST_BASE_URL"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
