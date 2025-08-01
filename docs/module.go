package docs

import (
	"context"
	"github.com/bd878/doc_server/internal/system"
	"github.com/bd878/doc_server/internal/grpc"
	"github.com/bd878/doc_server/docs/internal/controller"
	"github.com/bd878/doc_server/docs/internal/handlers"
	"github.com/bd878/doc_server/docs/internal/repository"
	"github.com/bd878/doc_server/docs/internal/gateway/users"
)

type Module struct {
}

func (Module) Startup(ctx context.Context, mono system.Service) (err error) {
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	gateway := users.NewGateway(conn)

	docs := repository.New("docs.meta", "docs.files", "docs.permissions", mono.DB())
	ctrl := controller.New(docs, gateway)

	handlers.RegisterHandlers(mono.Mux(), ctrl, mono.Logger())

	return nil
}
