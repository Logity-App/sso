package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
)

type Config struct {
	App  AppConfig
	GRPC GRPCCConfig
}

type AppConfig struct {
	Env         string        `env:"ENV" env-default:"local"`
	StoragePath string        `env:"STORAGE_PATH" env-required:"true"`
	TokenTTL    time.Duration `env:"TOKEN_TTL" env-required:"true"`
}

type GRPCCConfig struct {
	Port    int           `env:"PORT" env-required:"true"`
	Timeout time.Duration `env:"TIMEOUT" env-default:"5s"`
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
