package repository

import (
	"context"
	"database/sql"

	"errors"
	"fmt"

	"github.com/lib/pq"

	"subscriptions-api/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Insert(ctx context.Context, s models.Subscription) error {
	query := `
		INSERT INTO subscriptions 
		    (id, user_id, service, price, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		s.ID, s.UserID, s.Service, s.Price, s.StartDate, s.EndDate,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				return ErrDuplicate
			case "23502": // not_null_violation
				return fmt.Errorf("%w: missing required field %s", ErrDB, pqErr.Column)
			default:
				return fmt.Errorf("%w: %v", ErrDB, err)
			}
		}
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	return nil
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	query := `
		SELECT id, user_id, service, price, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var s models.Subscription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.Service,
		&s.Price,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Subscription{}, ErrNotFound
		}
		return models.Subscription{}, fmt.Errorf("%w: %v", ErrDB, err)
	}

	return s, nil
}

func (r *SubscriptionRepo) List(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, user_id, service, price, start_date, end_date, created_at, updated_at
		FROM subscriptions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]models.Subscription, 0)
	for rows.Next() {
		var sub models.Subscription
		var endDate sql.NullTime
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.Service, &sub.Price, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, err
		}
		if endDate.Valid {
			sub.EndDate = &endDate.Time
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service = $1,
		    price = $2,
		    start_date = $3,
		    end_date = $4,
		    updated_at = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(ctx, query,
		sub.Service,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.UpdatedAt,
		sub.ID,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return ErrDuplicate
			case "23502":
				return fmt.Errorf("%w: missing required field %s", ErrDB, pqErr.Column)
			default:
				return fmt.Errorf("%w: %v", ErrDB, err)
			}
		}
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		// Общая ошибка базы
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SubscriptionRepo) Sum(ctx context.Context, userID *uuid.UUID, service *string, start, end string) (int, error) {
	query := `SELECT COALESCE(SUM(price),0) FROM subscriptions WHERE start_date >= $1 AND start_date <= $2`
	args := []interface{}{start, end}

	// фильтр по user_id
	if userID != nil {
		query += ` AND user_id = $3`
		args = append(args, *userID)
	}

	// фильтр по service
	if service != nil {
		if userID != nil {
			query += ` AND service = $4`
		} else {
			query += ` AND service = $3`
		}
		args = append(args, *service)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("db query error: %w", err)
	}

	return total, nil
}
