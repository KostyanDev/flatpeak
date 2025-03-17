package config

import (
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

func New[T any]() (T, error) {
	var cfg T

	if envFilePath, ok := os.LookupEnv("ENV_FILE_PATH"); ok {
		if err := godotenv.Load(envFilePath); err != nil {
			return cfg, err
		}
	}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	if ok, err := govalidator.ValidateStruct(cfg); !ok {
		return cfg, err
	}

	return cfg, nil
}
