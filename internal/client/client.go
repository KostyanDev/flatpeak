package client

import (
	"app/internal/config"
	"app/internal/domain"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.30.0 --name=CarbonClient --output=internal/mocks --case=underscore
type CarbonClient interface {
	FetchCarbonForecast(ctx context.Context) ([]domain.Carbon, error)
}

type Client struct {
	cfg        config.Config
	log        *logrus.Logger
	httpClient *http.Client
}

func NewClient(cfg config.Config, log *logrus.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cfg: cfg,
		log: log,
	}
}
