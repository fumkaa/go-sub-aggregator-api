package httpserver

import (
	"log/slog"
)

type SubscriptionService interface {
}

type serverAPI struct {
	subService SubscriptionService
	log        *slog.Logger
}

func NewServerAPI(subService SubscriptionService, log *slog.Logger) *serverAPI {
	return &serverAPI{
		subService: subService,
		log:        log,
	}
}
