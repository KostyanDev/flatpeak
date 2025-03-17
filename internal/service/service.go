package service

import (
	"app/internal/client"
	"context"

	"github.com/sirupsen/logrus"
)

type Service struct {
	ctx          context.Context
	log          *logrus.Logger
	carbonClient client.CarbonClient
}

func New(ctx context.Context, log *logrus.Logger, carbonClient client.CarbonClient) *Service {
	return &Service{
		ctx:          ctx,
		log:          log,
		carbonClient: carbonClient,
	}
}
