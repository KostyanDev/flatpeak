package http

import (
	"app/internal/domain"
	"app/internal/service"
	"context"

	"github.com/sirupsen/logrus"
)

type Service interface {
	GetWeightedCarbonIntensity(ctx context.Context, filter domain.GetSlots) ([]domain.LowestCarbonPeriod, error)
}

type Handler struct {
	ctx     context.Context
	log     *logrus.Logger
	service Service
}

func New(ctx context.Context, log *logrus.Logger, service *service.Service) *Handler {
	return &Handler{
		ctx:     ctx,
		log:     log,
		service: service,
	}
}
