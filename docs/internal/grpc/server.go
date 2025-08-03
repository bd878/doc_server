package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/bd878/doc_server/docs/docspb"
)

type Controller interface {
	FreeCache(ctx context.Context, login string) (err error)
}

type server struct {
	ctrl Controller
	docspb.UnimplementedDocsServiceServer
}

var _ docspb.DocsServiceServer = (*server)(nil)

func RegisterServer(ctrl Controller, registrar *grpc.Server) {
	docspb.RegisterDocsServiceServer(registrar, server{ctrl: ctrl})
}

func (s server) FreeMemory(ctx context.Context, request *docspb.FreeMemoryRequest) (
	*docspb.FreeMemoryResponse, error,
) {
	err := s.ctrl.FreeCache(ctx, request.Login)
	if err != nil {
		return nil, err
	}

	return &docspb.FreeMemoryResponse{}, nil
}
