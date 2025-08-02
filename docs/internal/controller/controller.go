package controller

import (
	"context"
	"io"
	"mime/multipart"
	"encoding/json"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Repository interface {
	SaveFile(ctx context.Context, f multipart.File, meta *docs.Meta) (err error)
	SaveJSON(ctx context.Context, data []byte, meta *docs.Meta) (err error)
	List(ctx context.Context, key, value string, limit int) (docs []*docs.Meta, isLastPage bool, err error)
	GetMeta(ctx context.Context, id string) (meta *docs.Meta, err error)
	ReadFile(ctx context.Context, id string) (file io.Reader, err error)
	ReadJSON(ctx context.Context, id string) (json json.RawMessage, err error)
	Delete(ctx context.Context, id string) (err error)
}

type Controller struct {
	repo     Repository
}

func New(repo Repository) *Controller {
	return &Controller{repo}
}

func (c Controller) SaveFile(ctx context.Context, f multipart.File, meta *docs.Meta) (err error) {
	return
}

func (c Controller) SaveJSON(ctx context.Context, json []byte, meta *docs.Meta) (err error) {
	return
}

func (c Controller) List(ctx context.Context, key, value string, limit int) (docs []*docs.Meta, isLastPage bool, err error) {
	return
}

func (c Controller) GetMeta(ctx context.Context, id string) (doc *docs.Meta, err error) {
	return
}

func (c Controller) ReadFileStream(ctx context.Context, id string) (file io.Reader, err error) {
	return
}

func (c Controller) ReadJSON(ctx context.Context, id string) (json json.RawMessage, err error) {
	return
}

func (c Controller) Delete(ctx context.Context, id string) (err error) {
	return
}
