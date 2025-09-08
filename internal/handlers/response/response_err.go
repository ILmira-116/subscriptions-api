package response

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error     string `json:"error" example:"invalid request body"`
	Code      string `json:"code,omitempty" example:"400"`
	RequestID string `json:"request_id,omitempty" example:"req-12345"`
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}

type ValidationErrorResponse struct {
	Error  string            `json:"error" example:"validation failed"`
	Fields map[string]string `json:"fields" example:"{\"UserID\":\"required\",\"Price\":\"must be >0\"}"`
}

func ValidationError(errs validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for _, fe := range errs {
		result[fe.Field()] = fe.Error()
	}
	return result
}

// Для swagger-документации
// Error400 пример для 400 Bad Request
type Error400 struct {
	Code      string `json:"code" example:"400"`
	Error     string `json:"error" example:"invalid request body"`
	RequestID string `json:"request_id,omitempty" example:"req-12345"`
}

// Error409 пример для 409 Conflict
type Error409 struct {
	Code      string `json:"code" example:"409"`
	Error     string `json:"error" example:"subscription already exists"`
	RequestID string `json:"request_id,omitempty" example:"req-45678"`
}

// Error500 пример для 500 Internal Server Error
type Error500 struct {
	Code      string `json:"code" example:"500"`
	Error     string `json:"error" example:"internal server error"`
	RequestID string `json:"request_id,omitempty" example:"req-78901"`
}

// Error404 пример для 404 Not Found
type Error404 struct {
	Code      string `json:"code" example:"404"`
	Error     string `json:"error" example:"subscription not found"`
	RequestID string `json:"request_id,omitempty" example:"req-99999"`
}
