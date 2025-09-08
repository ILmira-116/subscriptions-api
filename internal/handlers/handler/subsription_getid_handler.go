package handler

import (
	"context"
	"net/http"
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

// GetSubscription godoc
// @Summary Получить подписку по ID
// @Description Возвращает данные подписки по её UUID
// @Tags subscriptions
// @Produce  json
// @Param   id path string true "ID подписки (UUID)"
// @Success 200 {object} response.GetSubscriptionResponse "Подписка найдена"
// @Failure 400 {object} response.Error400 "Неверный UUID"
// @Failure 404 {object} response.Error404 "Подписка не найдена"
// @Failure 500 {object} response.Error500 "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func NewGetSubscriptionHandler(log *logger.Logger, svc SubscriptionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID из URL
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid UUID"))
			return
		}

		// Вызываем сервис
		sub, err := svc.GetSubscription(r.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrSubscriptionNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("subscription not found"))
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
