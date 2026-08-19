package auth_service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	core_errors "github.com/c-dvoid/teamwork_api_go/internal/core/errors"
)

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return "", fmt.Errorf("login user: %w", core_errors.ErrUnauthorized)
		}
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("login user: %w", core_errors.ErrUnauthorized)
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}
