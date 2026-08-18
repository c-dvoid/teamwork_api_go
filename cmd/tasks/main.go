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
	core_http_router "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/router.go"
	core_http_server "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/server"
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

	r := core_http_router.NewRouter(log)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	server := core_http_server.NewServer(httpCfg, log, r)

	if err := server.Run(ctx); err != nil {
		log.Error("server run failed", zap.Error(err))
	}
}
