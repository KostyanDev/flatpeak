package config

import (
	"time"
)

type Config struct {
	App          App
	HTTPServer   HTTPServer   `envPrefix:"HTTP_"`
	CarbonClient CarbonClient `envPrefix:"CARBON_CLIENT"`
}

type App struct {
	Name string `env:"APP_NAME" envDefault:"app"`
}

type HTTPServer struct {
	Host         string        `env:"SERVER_HOST" envDefault:"localhost"`
	Port         int           `env:"SERVER_PORT" envDefault:"8080"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"15s"`
}

type CarbonClient struct {
	URL string `env:"URL" envDefault:"https://api.carbonintensity.org.uk"`
}
