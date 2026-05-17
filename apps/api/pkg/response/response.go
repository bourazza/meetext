package response

import (
	"encoding/json"
	"net/http"

	"github.com/meetext/backend/pkg/apperr"
)

type envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *errBody    `json:"error,omitempty"`
}

type errBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Success: true, Data: data})
}

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	appErr, ok := apperr.As(err)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(envelope{
			Success: false,
			Error:   &errBody{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	w.WriteHeader(appErr.StatusCode)
	_ = json.NewEncoder(w).Encode(envelope{
		Success: false,
		Error:   &errBody{Code: appErr.Code, Message: appErr.Message},
	})
}

func ValidationError(w http.ResponseWriter, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(envelope{
		Success: false,
		Error: &errBody{
			Code:    "VALIDATION_ERROR",
			Message: "validation failed",
			Fields:  fields,
		},
	})
}
