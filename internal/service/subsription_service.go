package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"subscriptions-api/internal/handlers/request/payload"
	"subscriptions-api/internal/models"
	"subscriptions-api/internal/repository"
	"subscriptions-api/pkg/utils/logger"

	"github.com/google/uuid"
)

type SubscriptionRepoInsert interface {
	Insert(ctx context.Context, s models.Subscription) error
}

type SubscriptionRepoGet interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.Subscription, error)
}

type SubscriptionRepoList interface {
	List(ctx context.Context, limit, offset int) ([]models.Subscription, error)
}

type SubscriptionRepoUpdate interface {
	Update(ctx context.Context, s models.Subscription) error
}

type SubscriptionRepoDelete interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type SubscriptionSummarizer interface {
	Sum(ctx context.Context, userID *uuid.UUID, service *string, start, end string) (int, error)
}

type SubscriptionService struct {
	repo interface {
		SubscriptionRepoInsert
		SubscriptionRepoGet
		SubscriptionRepoList
		SubscriptionRepoUpdate
		SubscriptionRepoDelete
		SubscriptionSummarizer
	}
	log *logger.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepo, log *logger.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
		log:  log,
	}
}
func (s *SubscriptionService) CreateSubscription(ctx context.Context, payload payload.CreateSubscriptionPayload) (uuid.UUID, error) {
	// Конвертируем StartDate
	start, endDefault, err := parseMonthYear(payload.StartDate)
	if err != nil {
		s.log.Error("invalid start_date format", "start_date", payload.StartDate, "error", err)
		return uuid.Nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	var end time.Time
	if payload.EndDate != nil {
		_, endParsed, err := parseMonthYear(*payload.EndDate)
		if err != nil {
			s.log.Error("invalid end_date format", "end_date", *payload.EndDate, "error", err)
			return uuid.Nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		end = endParsed
	} else {
		end = endDefault
	}

	sub := models.Subscription{
		ID:        uuid.New(),
		UserID:    payload.UserID,
		Service:   payload.Service,
		Price:     payload.Price,
		StartDate: start,
		EndDate:   &end,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Вставка в репозиторий
	err = s.repo.Insert(ctx, sub)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			s.log.Warn("subscription already exists",
				"user_id", sub.UserID.String(),
				"service", sub.Service,
			)
			return uuid.Nil, ErrSubscriptionAlreadyExists
		}

		s.log.Error("failed to insert subscription",
			"error", err,
			"user_id", sub.UserID.String(),
		)
		return uuid.Nil, err
	}

	s.log.Info("subscription created successfully",
		"subscription_id", sub.ID.String(),
		"user_id", sub.UserID.String(),
	)

	return sub.ID, nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.log.Warn("subscription not found",
				"subscription_id", id.String(),
			)
			return models.Subscription{}, ErrSubscriptionNotFound
		}

		s.log.Error("failed to get subscription",
			"error", err,
			"subscription_id", id.String(),
		)
		return models.Subscription{}, err
	}

	s.log.Info("subscription retrieved successfully",
		"subscription_id", id.String(),
		"user_id", sub.UserID.String(),
	)

	return sub, nil
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, limit, offset int) ([]models.Subscription, error) {

	subs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.log.Error("failed to list subscriptions",
			"error", err,
			"limit", limit,
			"offset", offset,
		)
		return nil, err
	}

	s.log.Info("subscriptions retrieved successfully",
		"count", len(subs),
		"limit", limit,
		"offset", offset,
	)

	return subs, nil
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, payload payload.UpdateSubscriptionPayload) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.log.Warn("subscription not found", "subscription_id", id.String())
			return ErrSubscriptionNotFound
		}
		s.log.Error("failed to get subscription", "subscription_id", id.String(), "error", err)
		return err
	}

	// Конвертируем StartDate
	start, endDefault, err := parseMonthYear(payload.StartDate)
	if err != nil {
		s.log.Error("invalid start_date format", "start_date", payload.StartDate, "error", err)
		return fmt.Errorf("invalid start_date format: %w", err)
	}

	// Конвертируем EndDate, если передан
	var end time.Time
	if payload.EndDate != nil {
		_, endParsed, err := parseMonthYear(*payload.EndDate)
		if err != nil {
			s.log.Error("invalid end_date format", "end_date", *payload.EndDate, "error", err)
			return fmt.Errorf("invalid end_date format: %w", err)
		}
		end = endParsed
	} else {
		end = endDefault
	}

	// Обновляем поля подписки
	sub.Service = payload.Service
	sub.Price = payload.Price
	sub.StartDate = start
	sub.EndDate = &end
	sub.UpdatedAt = time.Now()

	// Сохраняем изменения
	err = s.repo.Update(ctx, sub)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicate):
			s.log.Warn("subscription already exists", "user_id", sub.UserID.String(), "service", sub.Service)
			return ErrSubscriptionAlreadyExists
		default:
			s.log.Error("failed to update subscription", "subscription_id", sub.ID.String(), "error", err)
			return err
		}
	}

	s.log.Info("subscription updated successfully", "subscription_id", sub.ID.String(), "user_id", sub.UserID.String())
	return nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			s.log.Warn("subscription not found",
				"subscription_id", id.String(),
			)
			return ErrSubscriptionNotFound

		default:
			s.log.Error("failed to delete subscription",
				"error", err,
				"subscription_id", id.String(),
			)
			return err
		}
	}

	s.log.Info("subscription deleted successfully",
		"subscription_id", id.String(),
	)

	return nil
}

func (s *SubscriptionService) SumSubscriptions(ctx context.Context, p payload.SubscriptionSummaryPayload) (int, error) {
	total, err := s.repo.Sum(ctx, p.UserID, p.Service, p.StartDate.Format("2006-01-02"), p.EndDate.Format("2006-01-02"))
	if err != nil {
		s.log.Error("failed to sum subscriptions",
			"error", err,
			"user_id", p.UserID,
			"service", p.Service,
			"start_date", p.StartDate,
			"end_date", p.EndDate,
		)
		return 0, err
	}

	s.log.Info("subscriptions summed successfully",
		"user_id", p.UserID,
		"service", p.Service,
		"start_date", p.StartDate,
		"end_date", p.EndDate,
		"total", total,
	)

	return total, nil
}

func parseMonthYear(dateStr string) (time.Time, time.Time, error) {
	t, err := time.Parse("01-2006", dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)

	return start, end, nil
}
