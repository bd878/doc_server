package repository

import (
	"context"
	"database/sql"
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

func (r *Repository) Save(ctx context.Context, doc *docs.Doc) (err error) {
	return
}

func (r *Repository) List(ctx context.Context, key, value string, limit int) (docs []*docs.Doc, isLastPage bool, err error) {
	return
}

func (r *Repository) Get(ctx context.Context, id int) (doc *docs.Doc, err error) {
	return
}

func (r *Repository) Delete(ctx context.Context, id int) (err error) {
	return
}
