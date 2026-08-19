package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	core_logger "github.com/c-dvoid/teamwork_api_go/internal/core/logger"
	core_mysql_conn "github.com/c-dvoid/teamwork_api_go/internal/core/repository/mysql/conn"
	core_redis_conn "github.com/c-dvoid/teamwork_api_go/internal/core/repository/redis/conn"
	core_http_config "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/config"
	core_http_router "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/router"
	core_http_server "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/server"
	auth_repository_mysql "github.com/c-dvoid/teamwork_api_go/internal/features/auth/repository/mysql"
	auth_service "github.com/c-dvoid/teamwork_api_go/internal/features/auth/service"
	auth_transport_http "github.com/c-dvoid/teamwork_api_go/internal/features/auth/transport/http"
)

func main() {
	loggerCfg := core_logger.NewConfigMust()
	mysqlCfg := core_mysql_conn.NewConfigMust()
	redisCfg := core_redis_conn.NewConfigMust()
	httpCfg := core_http_config.NewConfigMust()

	log, err := core_logger.NewLogger(loggerCfg)
	if err != nil {
		panic(err)
	}
	defer log.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mysqlPool := core_mysql_conn.NewPoolMust(ctx, mysqlCfg)
	defer mysqlPool.Close()
	log.Info("mysql connected")

	redisPool := core_redis_conn.NewPoolMust(ctx, redisCfg)
	defer redisPool.Close()
	log.Info("redis connected")

	authRepo := auth_repository_mysql.NewAuthRepository(mysqlPool.DB)
	authService := auth_service.NewAuthService(authRepo, httpCfg.JWTSecret, httpCfg.JWTTTL)
	authHandler := auth_transport_http.NewAuthHTTPHandler(authService)

	r := core_http_router.NewRouter(log)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	authHandler.RegisterRoutes(r)

	server := core_http_server.NewServer(httpCfg, log, r)

	if err := server.Run(ctx); err != nil {
		log.Error("server run failed", zap.Error(err))
	}
}
