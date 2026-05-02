package dto

import (
	"encoding/json"
	"net/http"
)

type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ApiErrorResponse[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ApiSuccessResponse[T any] struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Record fetched successfully"`
	Data    T      `json:"data"`
}

func JsonResponse[T any](w http.ResponseWriter, response *ApiResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}

func JsonErrorResponse[T any](w http.ResponseWriter, response *ApiResponse[T]) {
	JsonResponse(w, &ApiResponse[string]{
		Code:    response.Code,
		Status:  "error",
		Message: response.Message,
	})
}

func JsonNotFoundResponse(w http.ResponseWriter, message string) {
	JsonResponse(w, &ApiResponse[string]{
		Code:    http.StatusNotFound,
		Status:  http.StatusText(http.StatusNotFound),
		Message: message,
	})
}

func JsonBadRequestResponse(w http.ResponseWriter, message string) {
	JsonResponse(w, &ApiResponse[string]{
		Code:    http.StatusBadRequest,
		Status:  http.StatusText(http.StatusBadRequest),
		Message: message,
	})
}

func JsonInternalServerErrorResponse(w http.ResponseWriter, message string) {
	JsonResponse(w, &ApiResponse[string]{
		Code:    http.StatusInternalServerError,
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: message,
	})
}

func JsonSuccessResponse[T any](w http.ResponseWriter, response *ApiSuccessResponse[T]) {
	JsonResponse(w, &ApiResponse[T]{
		Code:    http.StatusOK,
		Status:  "success",
		Message: response.Message,
		Data:    response.Data,
	})
}
