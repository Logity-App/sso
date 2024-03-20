package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
)

type Config struct {
	App      AppConfig
	GRPC     GRPCCConfig
	Database Database
}

type AppConfig struct {
	Env      string        `env:"ENVIRONMENT" env-default:"local"`
	TokenTTL time.Duration `env:"TOKEN_TTL" env-required:"true"`
}

type GRPCCConfig struct {
	Port    int           `env:"PORT" env-required:"true"`
	Timeout time.Duration `env:"TIMEOUT" env-default:"5s"`
}

type Database struct {
	Enable          bool   `env:"POSTGRES_ENABLE" env-default:"false"`
	Dsn             string `env:"POSTGRES_DSN" env-default:"postgres:5432"`
	MaxIdleConn     int    `env:"POSTGRES_MAX_IDLE_CONN" env-default:"3"`
	MaxLifeTimeConn int    `env:"POSTGRES_LIFETIME_CONN" env-default:"3"`
}

var configInstance *Config
var configErr error

func GetConfig() (*Config, error) {
	if configInstance == nil {
		var readConfigOnce sync.Once

		readConfigOnce.Do(func() {
			configInstance = &Config{}
			configErr = cleanenv.ReadEnv(configInstance)
		})
	}

	return configInstance, configErr
}
