// Package config ...
package config

import (
	"github.com/caarlos0/env/v6"
)

// Config configuration
type Config struct {
	PostgresURL string `env:"POSTGRES_URL"`
	MongoURL    string `env:"MONGO_URL"`
	ServerPort  string `env:"SERVER_ADDRESS"`
	DBType      string `env:"DB_TYPE"`
	RedisURL    string `env:"REDIS_URL"`
}

// New configuration
func New() (*Config, error) {
	var config Config
	err := env.Parse(&config)

	return &config, err
}
