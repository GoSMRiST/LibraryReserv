package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	HostAddress string        `env:"REST_HOST_ADDRESS"`
	ServTimeout time.Duration `env:"SERV_TIMEOUT"`
	DBHost      string        `env:"DB_HOST"`
	DBPort      int           `env:"DB_PORT"`
	DBUser      string        `env:"DB_USER"`
	DBPassword  string        `env:"DB_PASSWORD"`
	DBName      string        `env:"DB_NAME"`
	LogLevel    string        `env:"LOG_LEVEL"`
}

func InitConfig() *Config {
	conf := Config{}

	err := godotenv.Load("internal/config/config.env")
	if err != nil {
		panic("failed to load config")
	}

	err = env.Parse(&conf)
	if err != nil {
		panic("failed to parse config: " + err.Error())
	}

	return &conf
}
