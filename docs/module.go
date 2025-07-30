package docs

import (
	"context"
	"github.com/bd878/doc_server/internal/system"
)

type Module struct {
}

func (Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return nil
}
