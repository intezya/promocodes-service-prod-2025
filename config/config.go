package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	ServerPort       string `env:"SERVER_PORT"`
	PostgresConn     string `env:"POSTGRES_CONN"`
	RedisHost        string `env:"REDIS_HOST"`
	RedisPort        string `env:"REDIS_PORT"`
	AntifraudAddress string `env:"ANTIFRAUD_ADDRESS"`
	RandomSecret     string `env:"RANDOM_SECRET"`
}

func New() *Config {
	cfg := &Config{}
	//pkglib.ConfigLoader.NewWithExtraPath(cfg, ".env")
	_ = cleanenv.ReadConfig(".env", cfg)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	return cfg
}
