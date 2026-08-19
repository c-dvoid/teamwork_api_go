package auth_service

import (
	"context"
	"time"

	domain_user "github.com/c-dvoid/teamwork_api_go/internal/domain/user"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, passwordHash, name string) (domain_user.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain_user.User, error)
}

type AuthService struct {
	authRepository AuthRepository
	jwtSecret      string
	jwtTTL         time.Duration
}

func NewAuthService(
	authRepository AuthRepository,
	jwtSecret string,
	jwtTTL time.Duration,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		jwtSecret:      jwtSecret,
		jwtTTL:         jwtTTL,
	}
}
