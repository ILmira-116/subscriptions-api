package response

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}

type ValidationErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields"`
}

func ValidationError(errs validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for _, fe := range errs {
		result[fe.Field()] = fe.Error() // или fe.Tag() для типа ошибки
	}
	return result
}
