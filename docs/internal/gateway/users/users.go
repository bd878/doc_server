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

func (g usersGateway) Auth(ctx context.Context, token string) (login string, err error) {
	resp, err := g.client.Auth(ctx, &userspb.AuthRequest{Token: token})
	if err != nil {
		return "", err
	}

	return resp.User.Login, nil
}