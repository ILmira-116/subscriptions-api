package handler

import (
	"context"
	"net/http"
	"subscriptions-api/internal/handlers/request/payload"
	"subscriptions-api/internal/handlers/response"
	"subscriptions-api/internal/models"
	"subscriptions-api/internal/service"
	"subscriptions-api/pkg/utils/logger"

	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type SubscriptionGetter interface {
	GetSubscription(ctx context.Context, id uuid.UUID) (models.Subscription, error)
}

func NewGetSubscriptionHandler(log *logger.Logger, svc SubscriptionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID из URL
		idStr := chi.URLParam(r, "id")
		subPayload := payload.GetSubscriptionPayload{}

		// Валидируем UUID
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationErrorResponse{
				Error:  "invalid request",
				Fields: map[string]string{"id": "must be a valid UUID"},
			})
			return
		}
		subPayload.ID = id

		// Вызываем сервис
		sub, err := svc.GetSubscription(r.Context(), subPayload.ID)
		if err != nil {
			if errors.Is(err, service.ErrSubscriptionNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.ErrorResponse{
					Error: "subscription not found",
					Code:  "NOT_FOUND",
				})
				return
			}
			log.Error("failed to get subscription", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		// Успешный ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.SubscriptionFetched(sub))
	}
}
