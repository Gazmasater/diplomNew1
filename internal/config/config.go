package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type RunConfig struct {
	Address          string `env:"RUN_ADDRESS"`
	DefaultDBConnStr string `env:"DATABASE_URI"`
	AccrualAddress   string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	RedisAddress     string `env:"REDIS_URI"`
}

const (
	defaultServer = ":8080"

	defaultConnStr = "host=localhost port=5432 user=lew password=qwert dbname=diplom sslmode=disable"

	defaultAccrualAddress = "http://localhost:8888"
	defaultRedisAddress   = "localhost:6379" //localhost
)

func Load() RunConfig {
	cfg := RunConfig{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.StringVar(&cfg.DefaultDBConnStr, "d", defaultConnStr, "connection string")
	flag.StringVar(&cfg.AccrualAddress, "r", defaultAccrualAddress, "accrual connection string")
	flag.StringVar(&cfg.RedisAddress, "m", defaultRedisAddress, "redis connection string")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}
	return cfg
}
