package handler

import (
	"context"
	"errors"
	"net/http"
	"subscriptions-api/internal/handlers/response"
	"subscriptions-api/internal/models"
	"subscriptions-api/internal/service"
	"subscriptions-api/pkg/utils/logger"
	utils "subscriptions-api/pkg/utils/paginator"

	"github.com/go-chi/render"
)

type SubscriptionLister interface {
	ListSubscriptions(ctx context.Context, limit, offset int) ([]models.Subscription, error)
}

// ListSubscriptions godoc
// @Summary Получить список подписок
// @Description Возвращает список подписок с пагинацией
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Лимит записей" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} response.ListSubscriptionsResponse "Список подписок"
// @Failure 400 {object} response.Error400 "Неверные параметры запроса"
// @Failure 404 {object} response.Error404 "Подписки не найдены"
// @Failure 500 {object} response.Error500 "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func NewListSubscriptionsHandler(log *logger.Logger, svc SubscriptionLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Парсим и валидируем пагинацию
		limit, offset, fields, err := utils.ParsePaginationParams(r)
		if err != nil {
			log.Warn("invalid pagination parameters", "fields", fields)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationErrorResponse{
				Error:  "invalid request",
				Fields: fields,
			})
			return
		}

		// Вызов сервиса
		subs, err := svc.ListSubscriptions(r.Context(), limit, offset)
		if err != nil {
			if errors.Is(err, service.ErrSubscriptionNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("no subscriptions found"))
				return
			}

			log.Error("failed to list subscriptions", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		// Формируем успешный ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.SubscriptionsFetched(subs))
	}

}
