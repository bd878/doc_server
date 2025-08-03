package controller

import (
	"context"
	"io"
	"mime/multipart"
	"encoding/json"
	"github.com/google/uuid"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, owner string, f multipart.File, json []byte, meta *docs.Meta) (err error)
	List(ctx context.Context, owner, login, key, value string, limit int) (docs []*docs.Meta, err error)
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

func (c Controller) Save(ctx context.Context, owner string, f multipart.File, json []byte, meta *docs.Meta) (err error) {
	meta.ID = uuid.New().String()

	err = c.repo.Save(ctx, owner, f, json, meta)

	return
}

func (c Controller) List(ctx context.Context, owner, login, key, value string, limit int) (docs []*docs.Meta, err error) {
	return c.repo.List(ctx, owner, login, key, value, limit)
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
