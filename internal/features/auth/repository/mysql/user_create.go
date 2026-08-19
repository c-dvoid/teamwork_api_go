package auth_repository_mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	core_errors "github.com/c-dvoid/teamwork_api_go/internal/core/errors"
	domain_user "github.com/c-dvoid/teamwork_api_go/internal/domain/user"
)

const mysqlDuplicateEntryCode = 1062

func (r *AuthRepository) CreateUser(ctx context.Context, email, passwordHash, name string) (domain_user.User, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, name) VALUES (?, ?, ?)`,
		email, passwordHash, name,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryCode {
			return domain_user.User{}, fmt.Errorf("create user: %w", core_errors.ErrAlreadyExists)
		}
		return domain_user.User{}, fmt.Errorf("insert user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain_user.User{}, fmt.Errorf("get last insert id: %w", err)
	}

	var row userRow
	err = r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&row.ID, &row.Email, &row.PasswordHash, &row.Name, &row.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain_user.User{}, fmt.Errorf("select created user: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain_user.User{}, fmt.Errorf("select created user: %w", err)
	}

	return domain_user.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		CreatedAt:    row.CreatedAt,
	}, nil
}
