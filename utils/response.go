package common_utils

import (
	"net/http"

	"github.com/kholiqcode/go-common/pkg/serializer"
)

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func GenerateSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if message == "" {
		message = "Success"
	}

	response := Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	responseEncode, err := serializer.Marshal(response)
	PanicIfAppError(err, "failed when marshal response", 500)

	if _, err := w.Write(responseEncode); err != nil {
		PanicIfAppError(err, "failed when write response", 500)
	}
}

func GenerateErrorResponse(w http.ResponseWriter, data interface{}, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if message == "" {
		message = "Something went wrong"
	}

	response := Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	responseEncode, err := serializer.Marshal(response)
	PanicIfAppError(err, "failed when marshar response", 500)

	if _, err := w.Write(responseEncode); err != nil {
		PanicIfAppError(err, "failed when write response", 500)
	}
}
