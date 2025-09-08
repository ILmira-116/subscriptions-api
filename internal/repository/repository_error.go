package repository

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("subscription not found")
	ErrDuplicate = errors.New("subscription already exists")
	ErrDB        = errors.New("database error")
)
