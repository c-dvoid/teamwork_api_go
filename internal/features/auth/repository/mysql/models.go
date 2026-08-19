package auth_repository_mysql

import "time"

type userRow struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
}
