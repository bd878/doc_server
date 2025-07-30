package system

import (
	"context"
	"net/http"
	"database/sql"

	"github.com/rs/zerolog"

	"github.com/bd878/doc_server/config"
	"github.com/bd878/doc_server/internal/waiter"
)

type Service interface {
	DB()      *sql.DB
	Mux()     *http.ServeMux
	Config()   config.AppConfig
	Waiter()  waiter.Waiter
	Logger()  zerolog.Logger
}

type Module interface {
	Startup(context.Context, Service) error
}