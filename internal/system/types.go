package system

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/doc_server/config"
	"github.com/bd878/doc_server/internal/waiter"
)

type Service interface {
	DB()      *pgxpool.Pool
	RPC()     *grpc.Server
	Mux()     *http.ServeMux
	Config()   config.AppConfig
	Waiter()   waiter.Waiter
	Logger()   zerolog.Logger
}

type Module interface {
	Startup(context.Context, Service) error
}