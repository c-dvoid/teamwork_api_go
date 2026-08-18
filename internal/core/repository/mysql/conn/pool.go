package core_mysql_conn

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 25
	defaultConnMaxLifetime = 5 * time.Minute
	defaultConnMaxIdleTime = 2 * time.Minute
)

type Pool struct {
	*sql.DB
	opTimeout time.Duration
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetConnMaxLifetime(defaultConnMaxLifetime)
	db.SetConnMaxIdleTime(defaultConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return &Pool{DB: db, opTimeout: config.Timeout}, nil
}

func NewPoolMust(ctx context.Context, config Config) *Pool {
	pool, err := NewPool(ctx, config)
	if err != nil {
		panic(fmt.Errorf("get MySQL connection pool: %w", err))
	}
	return pool
}
