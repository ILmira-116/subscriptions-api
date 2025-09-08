package handler

import (
	"context"
	"errors"
	"net/http"
	"subscriptions-api/internal/handlers/response"
	"subscriptions-api/internal/service"
	"subscriptions-api/pkg/utils/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SubscriptionDeleter interface {
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
}

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 200 {object} response.DeleteSubscriptionResponse "Подписка успешно удалена"
// @Failure 400 {object} response.Error400 "Неверный UUID"
// @Failure 404 {object} response.Error404 "Подписка не найдена"
// @Failure 500 {object} response.Error500 "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func NewDeleteSubscriptionHandler(log *logger.Logger, svc SubscriptionDeleter) http.HandlerFunc {
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

		// Вызываем сервис удаления
		if err := svc.DeleteSubscription(r.Context(), id); err != nil {
			switch {
			case errors.Is(err, service.ErrSubscriptionNotFound):
				errResp := response.Error("subscription not found")
				errResp.Code = "NOT_FOUND"
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, errResp)
				return
			default:
				log.Error("failed to delete subscription", "error", err, "subscription_id", id.String())
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal server error"))
				return
			}
		}

		log.Info("subscription deleted", "subscription_id", id.String())

		// Успешный ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.SubscriptionDeleted())
	}
}
