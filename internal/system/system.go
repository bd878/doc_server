package system

import (
	"fmt"
	"os"
	"context"
	"net/http"
	"database/sql"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"github.com/bd878/doc_server/config"
	"github.com/bd878/doc_server/internal/waiter"
	"github.com/bd878/doc_server/internal/logger"
)

type System struct {
	cfg      config.AppConfig
	db       *sql.DB
	modules  []Module
	logger   zerolog.Logger
	waiter   waiter.Waiter
	mux      *http.ServeMux
}

var _ Service = (*System)(nil)

func NewSystem(cfg config.AppConfig) (s *System, err error) {
	s = &System{cfg: cfg}

	if err = s.initDB(); err != nil {
		return nil, err
	}

	s.initMux()
	s.initWaiter()
	s.initLogger()

	return s, nil
}

func (s System) Config() config.AppConfig {
	return s.cfg
}

func (s *System) initDB() (err error) {
	s.db, err = sql.Open("pgx", s.cfg.PG.Conn)
	return err
}

func (s *System) initLogger() {
	s.logger = logger.New(logger.LogConfig{
		Environment: s.cfg.Environment,
		LogLevel: logger.Level(s.cfg.LogLevel),
	})
}

func (s *System) initMux() {
	s.mux = http.NewServeMux()
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) Logger() zerolog.Logger {
	return s.logger
}

func (s *System) DB() *sql.DB {
	return s.db
}

func (s *System) Mux() *http.ServeMux {
	return s.mux
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr: s.cfg.Web.Address(),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "web server started %s\n", s.Config().Web.Address())
		defer fmt.Fprintln(os.Stdout, "web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}
