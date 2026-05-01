package dto

import (
	"errors"
	"net/http"
)

type ApiError struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func (e *ApiError) ToResponse() *ApiResponse[any] {
	return &ApiResponse[any]{
		Code:    e.Code,
		Status:  http.StatusText(e.Code),
		Message: e.Message,
	}
}

func NewString(message string) *ApiError {
	return &ApiError{
		Code:    400,
		Message: message,
		Err:     errors.New(message),
	}
}

func CreateFromError(code int, err error) *ApiError {
	return &ApiError{
		Code:    code,
		Message: err.Error(),
		Err:     err,
	}
}
