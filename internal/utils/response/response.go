package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

var statusMessageMap = map[int]string{
	http.StatusOK:                  "ok",
	http.StatusCreated:             "created",
	http.StatusAccepted:            "accepted",
	http.StatusNoContent:           "no content",
	http.StatusBadRequest:          "bad request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusNotFound:            "not found",
	http.StatusConflict:            "conflict",
	http.StatusUnprocessableEntity: "unprocessable entity",
	http.StatusTooManyRequests:     "too many requests",
	http.StatusInternalServerError: "internal server error",
	http.StatusBadGateway:          "bad gateway",
	http.StatusServiceUnavailable:  "service unavailable",
	http.StatusGatewayTimeout:      "gateway timeout",
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf(err.Field()+" is required"))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Invalid value for field %s", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

func HandleError(w http.ResponseWriter, apiPath string, message string, statusCode int, err error) {
	errorMessage := message
	if errorMessage == "" {
		if err != nil {
			errorMessage = err.Error()
		} else {
			errorMessage = statusMessageByCode(statusCode)
		}
	}

	if err != nil {
		slog.Error(
			"API request failed",
			slog.String("path", apiPath),
			slog.Int("status_code", statusCode),
			slog.String("message", errorMessage),
			slog.String("error", err.Error()),
		)
	} else {
		slog.Error(
			"API request failed",
			slog.String("path", apiPath),
			slog.Int("status_code", statusCode),
			slog.String("message", errorMessage),
		)
	}

	_ = WriteJSON(w, statusCode, GeneralError(fmt.Errorf(errorMessage)))
}

func statusMessageByCode(statusCode int) string {
	if msg, ok := statusMessageMap[statusCode]; ok {
		return msg
	}

	if text := http.StatusText(statusCode); text != "" {
		return strings.ToLower(text)
	}

	return "request failed"
}
