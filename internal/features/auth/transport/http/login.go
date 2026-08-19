package auth_transport_http

import (
	"net/http"

	core_logger "github.com/c-dvoid/teamwork_api_go/internal/core/logger"
	core_http_request "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/request"
	core_http_response "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/response"
)

func (h *AuthHTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	resp := core_http_response.NewHTTPResponseHandler(log, w)

	var req LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(w, r, &req); err != nil {
		resp.ErrorResponse(err, "invalid request body")
		return
	}

	token, err := h.authService.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		resp.ErrorResponse(err, "failed to login")
		return
	}

	resp.JSONResponse(LoginResponse{Token: token}, http.StatusOK)
}
