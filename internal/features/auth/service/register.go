package auth_service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	domain_user "github.com/c-dvoid/teamwork_api_go/internal/domain/user"
)

func (s *AuthService) RegisterUser(ctx context.Context, email, password, name string) (domain_user.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain_user.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.authRepository.CreateUser(ctx, email, string(passwordHash), name)
	if err != nil {
		return domain_user.User{}, fmt.Errorf("register user: %w", err)
	}

	return user, nil
}
