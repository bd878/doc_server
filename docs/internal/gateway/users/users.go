package users

import (
	"context"
	"google.golang.org/grpc"
	"github.com/bd878/doc_server/users/userspb"
)

type usersGateway struct {
	client userspb.UsersServiceClient
}

func NewGateway(conn *grpc.ClientConn) *usersGateway {
	return &usersGateway{client: userspb.NewUsersServiceClient(conn)}
}

func (g usersGateway) Auth(ctx context.Context, token string) (ok bool, err error) {
	return
}