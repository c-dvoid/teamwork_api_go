package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/c-dvoid/teamwork_api_go/internal/core/errors"
	core_logger "github.com/c-dvoid/teamwork_api_go/internal/core/logger"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {

	if statusCode == http.StatusNoContent || responseBody == nil {
		h.rw.WriteHeader(statusCode)
		return
	}

	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}

}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrForbidden):
		statusCode = http.StatusForbidden
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrAlreadyExists):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Debug
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	errText := err.Error()

	// Если произойдет внутренняя ошибка сервера (500) - скрываем системный текст от клиента
	if statusCode == http.StatusInternalServerError {
		errText = "internal server error"
	}

	response := ErrorResponse{
		Message: msg,
		Error:   errText,
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
