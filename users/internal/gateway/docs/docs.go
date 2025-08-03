package docs

import (
	"context"
	"google.golang.org/grpc"
	"github.com/bd878/doc_server/docs/docspb"
)

type docsGateway struct {
	client docspb.DocsServiceClient
}

func NewGateway(conn *grpc.ClientConn) *docsGateway {
	return &docsGateway{client: docspb.NewDocsServiceClient(conn)}
}

func (g docsGateway) FreeMemory(ctx context.Context, login string) (err error) {
	_, err = g.client.FreeMemory(ctx, &docspb.FreeMemoryRequest{Login: login})

	return
}