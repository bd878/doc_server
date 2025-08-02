package repository

import (
	"context"
	"database/sql"
	"io"
	"encoding/json"
	"mime/multipart"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Repository struct {
	tableName              string
	metaTableName          string
	permissionsTableName   string
	pool                  *sql.DB
}

func New(tableName, metaTableName, permissionsTableName string, pool *sql.DB) *Repository {
	return &Repository{
		tableName:             tableName,
		metaTableName:         metaTableName,
		permissionsTableName:  permissionsTableName,
		pool:                  pool,
	}
}

func (r *Repository) SaveFile(ctx context.Context, f multipart.File, meta *docs.Meta) (err error) {
	return
}

func (r *Repository) SaveJSON(ctx context.Context, data []byte, meta *docs.Meta) (err error) {
	return
}

func (r *Repository) List(ctx context.Context, key, value string, limit int) (docs []*docs.Meta, isLastPage bool, err error) {
	return
}

func (r *Repository) GetMeta(ctx context.Context, id string) (meta *docs.Meta, err error) {
	return
}

func (r *Repository) ReadFile(ctx context.Context, id string) (file io.Reader, err error) {
	return
}

func (r *Repository) ReadJSON(ctx context.Context, id string) (json json.RawMessage, err error) {
	return
}

func (r *Repository) Delete(ctx context.Context, id string) (err error) {
	return
}
