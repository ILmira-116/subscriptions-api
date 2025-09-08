package response

import "subscriptions-api/internal/models"

// Ответ на создание записи
type CreateSubscriptionResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func SubscriptionCreated(id string) CreateSubscriptionResponse {
	return CreateSubscriptionResponse{
		ID:      id,
		Message: "subscription created successfully",
	}
}

// Ответ на получение записи
type GetSubscriptionResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Service   string  `json:"service"`
	Price     int     `json:"price"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
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

// Ответ на получение списка подписок
type ListSubscriptionsResponse struct {
	Subscriptions []GetSubscriptionResponse `json:"subscriptions"`
	Count         int                       `json:"count"`
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

// Ответ на обновление записи
type UpdateSubscriptionResponse struct {
	Message string `json:"message"`
}

func SubscriptionUpdated() UpdateSubscriptionResponse {
	return UpdateSubscriptionResponse{
		Message: "subscription updated successfully",
	}
}

// Ответ на удаление записи
type DeleteSubscriptionResponse struct {
	Message string `json:"message"`
}

func SubscriptionDeleted() DeleteSubscriptionResponse {
	return DeleteSubscriptionResponse{
		Message: "subscription deleted successfully",
	}
}

// Ответ на подсчёт суммарной стоимости подписок
type SubscriptionSummaryResponse struct {
	Total int `json:"total"`
}

func SubscriptionSummary(total int) SubscriptionSummaryResponse {
	return SubscriptionSummaryResponse{
		Total: total,
	}
}
