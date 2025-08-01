package grpc

import (
	"os"
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(ctx context.Context, endpoint string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err = conn.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing grpc conn: %w", err)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if err = conn.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing grpc conn: %w", err)
			}
		}()
	}()

	return
}
