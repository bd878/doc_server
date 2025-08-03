package main

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/bd878/doc_server/users"
	"github.com/bd878/doc_server/docs"
	"github.com/bd878/doc_server/config"
	"github.com/bd878/doc_server/internal/system"
)

type monolith struct {
	*system.System
	modules []system.Module
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "doc_server exited abnormally: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	m := &monolith{
		System: s,
		modules: []system.Module{
			&users.Module{},
			&docs.Module{},
		},
	}

	if err = m.startupModules(); err != nil {
		return err
	}

	m.Waiter().Add(
		m.WaitForWeb,
		m.WaitForRPC,
		m.WaitForDB,
	)

	return m.Waiter().Wait()
}

func (m *monolith) startupModules() error {
	for _, module := range m.modules {
		if err := module.Startup(m.Waiter().Context(), m); err != nil {
			return err
		}
	}
	return nil
}
