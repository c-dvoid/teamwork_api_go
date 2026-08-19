package core_http_router

import (
	"github.com/go-chi/chi/v5"

	core_logger "github.com/c-dvoid/teamwork_api_go/internal/core/logger"
	core_http_middleware "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/middleware"
)

func NewRouter(log *core_logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(core_http_middleware.RequestID())
	r.Use(core_http_middleware.Logger(log))
	r.Use(core_http_middleware.Panic())
	r.Use(core_http_middleware.Trace())

	return r
}
