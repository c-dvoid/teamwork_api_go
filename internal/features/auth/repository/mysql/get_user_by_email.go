package auth_repository_mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	core_errors "github.com/c-dvoid/teamwork_api_go/internal/core/errors"
	domain_user "github.com/c-dvoid/teamwork_api_go/internal/domain/user"
)

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (domain_user.User, error) {
	var row userRow

	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&row.ID, &row.Email, &row.PasswordHash, &row.Name, &row.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain_user.User{}, fmt.Errorf("get user by email: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain_user.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return domain_user.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		CreatedAt:    row.CreatedAt,
	}, nil
}
