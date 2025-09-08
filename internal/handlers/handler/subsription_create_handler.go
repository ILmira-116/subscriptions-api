package handler

import (
	"context"
	"net/http"
	"subscriptions-api/internal/handlers/request/payload"
	"subscriptions-api/internal/handlers/response"
	"subscriptions-api/internal/service"
	"subscriptions-api/pkg/utils/logger"

	"errors"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionCreator interface {
	CreateSubscription(ctx context.Context, payload payload.CreateSubscriptionPayload) (uuid.UUID, error)
}

func NewCreateSubscriptionHandler(log *logger.Logger, svc SubscriptionCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req payload.CreateSubscriptionPayload

		// Декодируем JSON
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("invalid request body", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		log.Info("request body decoded", "req", req)

		// Валидируем payload
		if err := req.Validate(); err != nil {
			log.Error("validation failed", "error", err)
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.ValidationError(ve))
				return
			}
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		// Вызываем сервис
		id, err := svc.CreateSubscription(r.Context(), req)
		if err != nil {
			if errors.Is(err, service.ErrSubscriptionAlreadyExists) {
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, response.Error("subscription already exists"))
				return
			}
			log.Error("failed to create subscription", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		log.Info("subscription created",
			"subscription_id", id.String(),
			"user_id", req.UserID.String(),
		)

		// Отправляем успешный ответ
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.SubscriptionCreated(id.String()))
	}
}
