package response

import "subscriptions-api/internal/models"

// Ответ на создание записи
type CreateSubscriptionResponse struct {
	ID      string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Message string `json:"message" example:"subscription created successfully"`
}

func SubscriptionCreated(id string) CreateSubscriptionResponse {
	return CreateSubscriptionResponse{
		ID:      id,
		Message: "subscription created successfully",
	}
}

// GetSubscriptionResponse — ответ на получение подписки
type GetSubscriptionResponse struct {
	ID        string  `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	UserID    string  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Service   string  `json:"service" example:"Yandex Plus"`
	Price     int     `json:"price" example:"400"`
	StartDate string  `json:"start_date" example:"07-2025"`
	EndDate   *string `json:"end_date,omitempty" example:"08-2025"`
	CreatedAt string  `json:"created_at" example:"2025-07-01 12:00:00"`
	UpdatedAt string  `json:"updated_at" example:"2025-07-01 12:00:00"`
}

func SubscriptionFetched(sub models.Subscription) GetSubscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		str := sub.EndDate.Format("2006-01-02")
		endDate = &str
	}

	return GetSubscriptionResponse{
		ID:        sub.ID.String(),
		UserID:    sub.UserID.String(),
		Service:   sub.Service,
		Price:     sub.Price,
		StartDate: sub.StartDate.Format("2006-01-02"),
		EndDate:   endDate,
		CreatedAt: sub.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: sub.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ListSubscriptionsResponse - Ответ на получение списка подписок
type ListSubscriptionsResponse struct {
	Subscriptions []GetSubscriptionResponse `json:"subscriptions"`
	Count         int                       `json:"count" example:"2"`
}

func SubscriptionsFetched(subs []models.Subscription) ListSubscriptionsResponse {
	result := make([]GetSubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		result = append(result, SubscriptionFetched(sub))
	}

	return ListSubscriptionsResponse{
		Subscriptions: result,
		Count:         len(result),
	}
}

// SubscriptionUpdated - Ответ на обновление записи
type SubscriptionUpdated struct {
	Message string `json:"message" example:"subscription updated successfully"`
}

func NewSubscriptionUpdated() SubscriptionUpdated {
	return SubscriptionUpdated{
		Message: "subscription updated successfully",
	}
}

//  DeleteSubscriptionResponse - Ответ на удаление записи
type DeleteSubscriptionResponse struct {
	Message string `json:"message" example:"subscription deleted successfully"`
}

func SubscriptionDeleted() DeleteSubscriptionResponse {
	return DeleteSubscriptionResponse{
		Message: "subscription deleted successfully",
	}
}

// SubscriptionSummaryResponse - Ответ на подсчёт суммарной стоимости подписок
type SubscriptionSummaryResponse struct {
	Total int `json:"total" example:"1200"`
}

func SubscriptionSummary(total int) SubscriptionSummaryResponse {
	return SubscriptionSummaryResponse{
		Total: total,
	}
}
