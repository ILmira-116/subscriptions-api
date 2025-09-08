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
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Service   string    `json:"service" validate:"required"`
	Price     int       `json:"price" validate:"required,gt=0"`
	StartDate string    `json:"start_date" validate:"required"`
	EndDate   *string   `json:"end_date,omitempty"`
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

// GetSubscriptionPayload — данные для запроса подписки по ID
type GetSubscriptionPayload struct {
	ID uuid.UUID `json:"-"`
}

// ParseAndValidate — конвертируем строку в UUID и проверяем
func ParseAndValidate(idStr string) (*GetSubscriptionPayload, error) {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("id must be a valid UUID")
	}

	return &GetSubscriptionPayload{
		ID: uid,
	}, nil
}

// UpdateSubscriptionPayload — данные для обновления подписки
type UpdateSubscriptionPayload struct {
	Service   string  `json:"service" validate:"required"`
	Price     int     `json:"price" validate:"required,gt=0"`
	StartDate string  `json:"start_date" validate:"required"`
	EndDate   *string `json:"end_date,omitempty"`
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
