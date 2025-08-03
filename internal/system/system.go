package system

import (
	"fmt"
	"os"
	"net"
	"time"
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/doc_server/config"
	"github.com/bd878/doc_server/internal/waiter"
	"github.com/bd878/doc_server/internal/logger"
)

type System struct {
	cfg      config.AppConfig
	db       *pgxpool.Pool
	modules  []Module
	logger   zerolog.Logger
	waiter   waiter.Waiter
	mux      *http.ServeMux
	rpc      *grpc.Server
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
	s.initRpc()

	return s, nil
}

func (s System) Config() config.AppConfig {
	return s.cfg
}

func (s *System) initDB() (err error) {
	s.db, err = pgxpool.New(context.Background(), s.cfg.PG.Conn)

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

func (s *System) initRpc() {
	server := grpc.NewServer()
	reflection.Register(server)

	s.rpc = server
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

func (s *System) DB() *pgxpool.Pool {
	return s.db
}

func (s *System) Mux() *http.ServeMux {
	return s.mux
}

func (s *System) RPC() *grpc.Server {
	return s.rpc
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
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

func (a *System) WaitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.cfg.Rpc.Address())
	if err != nil {
		return nil
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "rpc server started %s\n", a.Config().Rpc.Address())
		defer fmt.Fprintln(os.Stdout, "rpc server shutdown")
		if err := a.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			a.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(a.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			a.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (a *System) WaitForDB(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "closing pgpool connections")
		a.db.Close()
		return nil
	})

	return group.Wait()
}