package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/bd878/doc_server/users/userspb"
	"github.com/bd878/doc_server/users/pkg/model"
)

type Controller interface {
	Auth(ctx context.Context, token string) (user *model.User, err error)
}

type server struct {
	ctrl Controller
	userspb.UnimplementedUsersServiceServer
}

var _ userspb.UsersServiceServer = (*server)(nil)

func RegisterServer(ctrl Controller, registrar *grpc.Server) {
	userspb.RegisterUsersServiceServer(registrar, server{ctrl: ctrl})
}

func (s server) Auth(ctx context.Context, request *userspb.AuthRequest) (
	*userspb.AuthResponse, error,
) {
	user, err := s.ctrl.Auth(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	return &userspb.AuthResponse{
		User: &userspb.User{
			Login: user.Login,
		},
	}, nil
}
