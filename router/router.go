package router

import (
	handler "subscriptions-api/internal/handlers/handler"
	"subscriptions-api/pkg/utils/logger"

	"github.com/go-chi/chi/v5"
)

type AllSubscriptionServices interface {
	handler.SubscriptionCreator
	handler.SubscriptionGetter
	handler.SubscriptionLister
	handler.SubscriptionUpdater
	handler.SubscriptionDeleter
	handler.SubscriptionSummarizer
}

// NewRouter создает роутер с передачей объединенного сервиса
func NewRouter(log *logger.Logger, svc AllSubscriptionServices) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/subscriptions", handler.NewCreateSubscriptionHandler(log, svc))
	r.Get("/subscriptions/{id}", handler.NewGetSubscriptionHandler(log, svc))
	r.Get("/subscriptions", handler.NewListSubscriptionsHandler(log, svc))
	r.Put("/subscriptions/{id}", handler.NewUpdateSubscriptionHandler(log, svc, svc))
	r.Delete("/subscriptions/{id}", handler.NewDeleteSubscriptionHandler(log, svc))
	r.Get("/subscriptions/summary", handler.NewSubscriptionSummaryHandler(log, svc))

	return r
}
