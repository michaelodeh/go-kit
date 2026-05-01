package validate

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/michaelodeh/go-kit/dto"
)

func GetValidationMessage(err error, fallbackMessage string) *dto.ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		messages := make([]string, len(ve))
		for i, fe := range ve {
			messages[i] = fmt.Sprintf("%s is invalid", fe.Field())
		}
		return &dto.ApiError{
			Message: strings.Join(messages, ", "),
			Status:  http.StatusText(http.StatusBadRequest),
			Code:    http.StatusBadRequest,
		}
	}

	return &dto.ApiError{
		Message: fallbackMessage,
		Status:  http.StatusText(http.StatusBadRequest),
		Code:    http.StatusBadRequest,
	}
}
