package auth_transport_http

import (
	"context"

	domain_user "github.com/c-dvoid/teamwork_api_go/internal/domain/user"
	auth_service "github.com/c-dvoid/teamwork_api_go/internal/features/auth/service"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password, name string) (domain_user.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	ParseJWT(tokenString string) (*auth_service.Claims, error)
}

type AuthHTTPHandler struct {
	authService AuthService
}

func NewAuthHTTPHandler(
	authService AuthService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}
