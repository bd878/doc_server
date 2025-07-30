package main

import (
	"fmt"
	"os"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	cfg := config.LoadConfig(os.Args[1])

	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	m := &monolith{
		System: s,
		modules: []system.Module{
		},
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.DB())

	if err = m.startupModules(); err != nil {
		return err
	}

	m.Waiter().Add(
		m.WaitForWeb,
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
