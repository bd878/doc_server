package users

import (
	"context"
	"github.com/bd878/doc_server/internal/system"
	"github.com/bd878/doc_server/users/internal/controller"
	"github.com/bd878/doc_server/users/internal/handlers"
	"github.com/bd878/doc_server/users/internal/repository"
)

type Module struct {
}

func (Module) Startup(ctx context.Context, mono system.Service) (err error) {
	users := repository.New("users.users", mono.DB())
	ctrl := controller.New(users, mono.Config().AdminToken)

	handlers.RegisterHandlers(mono.Mux(), ctrl)

	return nil
}
