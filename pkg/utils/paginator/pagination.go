package utils

import (
	"errors"
	"net/http"
	"strconv"
)

// ParsePaginationParams парсит и валидирует query-параметры limit и offset
func ParsePaginationParams(r *http.Request) (limit, offset int, fields map[string]string, err error) {
	fields = make(map[string]string)
	limit = 100 // дефолтное значение
	offset = 0  // дефолтное значение

	if l := r.URL.Query().Get("limit"); l != "" {
		v, convErr := strconv.Atoi(l)
		if convErr != nil || v <= 0 {
			fields["limit"] = "must be a positive integer"
		} else {
			limit = v
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		v, convErr := strconv.Atoi(o)
		if convErr != nil || v < 0 {
			fields["offset"] = "must be a non-negative integer"
		} else {
			offset = v
		}
	}

	if len(fields) > 0 {
		err = errors.New("invalid pagination parameters")
	}

	return
}
