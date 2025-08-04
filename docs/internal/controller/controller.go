package controller

import (
	"context"
	"io"
	"strconv"
	"mime/multipart"
	"encoding/json"
	"github.com/google/uuid"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, owner string, f multipart.File, json []byte, meta *docs.Meta) (err error)
	List(ctx context.Context, owner, login, key string, value interface{}, limit int) (docs []*docs.Meta, err error)
	GetMeta(ctx context.Context, id, login string) (meta *docs.Meta, err error)
	ReadFile(ctx context.Context, oid uint32, writer io.Writer) (err error)
	ReadJSON(ctx context.Context, id string) (json json.RawMessage, err error)
	Delete(ctx context.Context, id string) (err error)
}

type Cache interface {
	Set(owner string, meta *docs.Meta)
	Get(id, login string) (meta *docs.Meta)
	List(login, key string, value interface{}, limit int) (docs []*docs.Meta)
	Free(login string)
	Remove(id string)
}

type Controller struct {
	repo     Repository
	cache    Cache
}

func New(repo Repository, cache Cache) *Controller {
	return &Controller{repo, cache}
}

func (c Controller) Save(ctx context.Context, owner string, f multipart.File, json []byte, meta *docs.Meta) (err error) {
	meta.ID = uuid.New().String()

	err = c.repo.Save(ctx, owner, f, json, meta)
	if err != nil {
		return
	}

	c.cache.Set(owner, meta)

	return
}

func (c Controller) List(ctx context.Context, owner, login, key, value string, limit int) (docs []*docs.Meta, err error) {
	var val interface{}
	if key == "file" || key == "public" {
		file, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		val = file
	} else {
		val = value
	}

	if owner != "" {
		docs = c.cache.List(owner, key, val, limit)
	} else {
		docs = c.cache.List(login, key, val, limit)
	}
	if docs == nil {
		return c.repo.List(ctx, owner, login, key, val, limit)
	}
	return
}

func (c Controller) GetMeta(ctx context.Context, id, login string) (doc *docs.Meta, err error) {
	doc = c.cache.Get(id, login)

	if doc == nil {
		return c.repo.GetMeta(ctx, id, login)
	}
	return
}

func (c Controller) ReadFileStream(ctx context.Context, oid uint32, writer io.Writer) (err error) {
	return c.repo.ReadFile(ctx, oid, writer)
}

func (c Controller) ReadJSON(ctx context.Context, id string) (json json.RawMessage, err error) {
	return c.repo.ReadJSON(ctx, id)
}

func (c Controller) Delete(ctx context.Context, id string) (err error) {
	c.cache.Remove(id)
	return c.repo.Delete(ctx, id)
}

func (c Controller) FreeCache(ctx context.Context, login string) (err error) {
	c.cache.Free(login)
	return nil
}