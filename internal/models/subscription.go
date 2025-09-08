package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Service   string
	Price     int
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
