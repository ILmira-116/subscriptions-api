package service

import (
	"errors"
)

var (
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrInvalidPaginationParams   = errors.New("invalid pagination parameters")
)
