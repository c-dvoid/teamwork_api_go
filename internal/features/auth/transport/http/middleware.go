package auth_transport_http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	core_errors "github.com/c-dvoid/teamwork_api_go/internal/core/errors"
	core_logger "github.com/c-dvoid/teamwork_api_go/internal/core/logger"
	core_http_middleware "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/middleware"
	core_http_response "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/response"
)

type ctxKey int

const userIDCtxKey ctxKey = iota

const authHeaderPrefix = "Bearer "

func (h *AuthHTTPHandler) RequireAuth() core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			resp := core_http_response.NewHTTPResponseHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, authHeaderPrefix) {
				resp.ErrorResponse(
					fmt.Errorf("missing or malformed Authorization header: %w", core_errors.ErrUnauthorized),
					"unauthorized",
				)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, authHeaderPrefix)

			claims, err := h.authService.ParseJWT(tokenString)
			if err != nil {
				resp.ErrorResponse(
					fmt.Errorf("parse token: %w", core_errors.ErrUnauthorized),
					"unauthorized",
				)
				return
			}

			ctx := context.WithValue(r.Context(), userIDCtxKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIDCtxKey).(int64)
	if !ok {
		panic("no user_id in context: RequireAuth middleware not applied")
	}
	return userID
}
