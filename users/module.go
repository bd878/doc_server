package users

import (
	"context"
	"github.com/bd878/doc_server/internal/system"
	"github.com/bd878/doc_server/internal/grpc"
	usersGrpc "github.com/bd878/doc_server/users/internal/grpc"
	"github.com/bd878/doc_server/users/internal/gateway/docs"
	"github.com/bd878/doc_server/users/internal/controller"
	"github.com/bd878/doc_server/users/internal/handlers"
	"github.com/bd878/doc_server/users/internal/repository"
)

type Module struct {
}

func (Module) Startup(ctx context.Context, mono system.Service) (err error) {
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	users := repository.New("users.users", mono.DB())
	gateway := docs.NewGateway(conn)

	ctrl := controller.New(users, gateway, mono.Config().AdminToken)

	handlers.RegisterHandlers(mono.Mux(), ctrl, mono.Logger())
	usersGrpc.RegisterServer(ctrl, mono.RPC())

	return nil
}
