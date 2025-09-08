package payload

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

// CreateSubscriptionPayload — данные, которые приходят в запросе
type CreateSubscriptionPayload struct {
	UserID    uuid.UUID `json:"user_id" validate:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Service   string    `json:"service" validate:"required" example:"Yandex Plus"`
	Price     int       `json:"price" validate:"required,gt=0" example:"400"`
	StartDate string    `json:"start_date" validate:"required" example:"07-2025"`
	EndDate   *string   `json:"end_date,omitempty" example:"08-2025"`
}

// Валидация структуры запроса
func (p *CreateSubscriptionPayload) Validate() error {
	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.EndDate != nil && *p.EndDate == "" {
		return fmt.Errorf("field EndDate must be a valid datetime")
	}

	return nil
}

// ParseAndValidateUUID конвертирует строку в UUID и проверяет корректность
func ParseAndValidateUUID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("id must be a valid UUID")
	}
	return id, nil
}

// UpdateSubscriptionPayload — данные для обновления подписки
type UpdateSubscriptionPayload struct {
	Service   string  `json:"service" validate:"required" example:"Yandex Plus"`
	Price     int     `json:"price" validate:"required,gt=0" example:"400"`
	StartDate string  `json:"start_date" validate:"required" example:"07-2025"`
	EndDate   *string `json:"end_date,omitempty" example:"08-2025"`
}

func (p *UpdateSubscriptionPayload) Validate() error {
	if err := validate.Struct(p); err != nil {
		return err
	}

	// Проверка EndDate, если она передана
	if p.EndDate != nil && *p.EndDate == "" {
		return fmt.Errorf("field EndDate must be a valid MM-YYYY")
	}

	return nil
}

// Для подсчёта суммарной стоимости всех подписок
type SubscriptionSummaryPayload struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Service   *string    `json:"service,omitempty"`
	StartDate time.Time  `json:"start_date" validate:"required"`
	EndDate   time.Time  `json:"end_date" validate:"required,gtfield=StartDate"`
}

func (p *SubscriptionSummaryPayload) Validate() error {
	return validate.Struct(p)
}
