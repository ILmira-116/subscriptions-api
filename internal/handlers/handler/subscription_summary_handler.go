package handler

import (
	"context"
	"errors"
	"net/http"
	"subscriptions-api/internal/handlers/request/payload"
	"subscriptions-api/internal/handlers/response"
	"subscriptions-api/pkg/utils/logger"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionSummarizer interface {
	SumSubscriptions(ctx context.Context, p payload.SubscriptionSummaryPayload) (int, error)
}

func NewSubscriptionSummaryHandler(log *logger.Logger, svc SubscriptionSummarizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p payload.SubscriptionSummaryPayload
		q := r.URL.Query()

		// Парсим user_id (опционально)
		if uid := q.Get("user_id"); uid != "" {
			id, err := uuid.Parse(uid)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.ValidationErrorResponse{
					Error:  "invalid request",
					Fields: map[string]string{"user_id": "must be a valid UUID"},
				})
				return
			}
			p.UserID = &id
		}

		// Парсим service (опционально)
		if s := q.Get("service"); s != "" {
			p.Service = &s
		}

		// Парсим start_date
		startStr := q.Get("start_date")
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationErrorResponse{
				Error:  "invalid request",
				Fields: map[string]string{"start_date": "must be YYYY-MM-DD"},
			})
			return
		}
		p.StartDate = start

		// Парсим end_date
		endStr := q.Get("end_date")
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationErrorResponse{
				Error:  "invalid request",
				Fields: map[string]string{"end_date": "must be YYYY-MM-DD"},
			})
			return
		}
		p.EndDate = end

		// Валидация payload через validator
		if err := p.Validate(); err != nil {
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

		// Вызов сервиса для подсчета суммы
		total, err := svc.SumSubscriptions(r.Context(), p)
		if err != nil {
			log.Error("failed to sum subscriptions", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		// Успешный ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.SubscriptionSummary(total))
	}
}
