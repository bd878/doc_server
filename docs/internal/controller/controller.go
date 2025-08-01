package controller

import (
	"context"
	"mime/multipart"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, doc *docs.Doc) (err error)
	List(ctx context.Context, key, value string, limit int) (docs []*docs.Doc, isLastPage bool, err error)
	Get(ctx context.Context, id int) (doc *docs.Doc, err error)
	Delete(ctx context.Context, id int) (err error)
}

type UsersGateway interface {
	Auth(ctx context.Context, token string) (ok bool, err error)
}

type Controller struct {
	repo     Repository
	gateway  UsersGateway
}

func New(repo Repository, gateway UsersGateway) *Controller {
	return &Controller{repo, gateway}
}

func (c Controller) Save(ctx context.Context, f multipart.File, meta docs.Meta) (err error) {
	return
}

func (c Controller) List(ctx context.Context, key, value string, limit int) (docs []*docs.Doc, isLastPage bool, err error) {
	return
}

func (c Controller) Get(ctx context.Context, id int) (doc *docs.Doc, err error) {
	return
}

func (c Controller) Delete(ctx context.Context, id int) (err error) {
	return
}
