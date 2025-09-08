package handler

import (
	"context"
	"errors"
	"net/http"
	"subscriptions-api/internal/handlers/request/payload"
	"subscriptions-api/internal/handlers/response"
	"subscriptions-api/internal/service"
	"subscriptions-api/pkg/utils/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SubscriptionUpdater interface {
	UpdateSubscription(ctx context.Context, id uuid.UUID, payload payload.UpdateSubscriptionPayload) error
}

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет существующую запись о подписке по ID
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   id path string true "ID подписки (UUID)"
// @Param   subscription body payload.UpdateSubscriptionPayload true "Данные для обновления"
// @Success 200 {object} response.SubscriptionUpdated "Подписка успешно обновлена"
// @Failure 400 {object} response.Error400 "Неверный UUID или JSON"
// @Failure 404 {object} response.Error404 "Подписка не найдена"
// @Failure 409 {object} response.Error409 "Подписка уже существует"
// @Failure 500 {object} response.Error500 "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func NewUpdateSubscriptionHandler(log *logger.Logger, svcupd SubscriptionUpdater, svcget SubscriptionDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID из URL
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationErrorResponse{
				Error:  "invalid request",
				Fields: map[string]string{"id": "must be a valid UUID"},
			})
			return
		}

		var req payload.UpdateSubscriptionPayload

		// Декодируем JSON
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("invalid request body", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		// Валидируем payload
		if err := req.Validate(); err != nil {
			log.Error("validation failed", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationErrorResponse{
				Error:  "validation error",
				Fields: map[string]string{"body": err.Error()},
			})
			return
		}

		// Вызываем сервис с ID и payload
		if err := svcupd.UpdateSubscription(r.Context(), id, req); err != nil {
			switch {
			case errors.Is(err, service.ErrSubscriptionAlreadyExists):
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, response.Error("subscription already exists"))
				return
			case errors.Is(err, service.ErrSubscriptionNotFound):
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("subscription not found"))
				return
			default:
				log.Error("failed to update subscription", "error", err, "subscription_id", id.String())
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal server error"))
				return
			}
		}

		log.Info("subscription updated", "subscription_id", id.String())

		// Успешный ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.NewSubscriptionUpdated())
	}
}
