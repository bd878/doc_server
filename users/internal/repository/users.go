package repository

import (
	"context"
	"database/sql"
	"github.com/bd878/doc_server/users/pkg/model"
)

type Repository struct {
	tableName string
	pool      *sql.DB
}

func New(tableName string, pool *sql.DB) *Repository {
	return &Repository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r Repository) Save(ctx context.Context, user *model.User) (err error) {
	return nil
}

func (r Repository) Find(ctx context.Context, token string) (user *model.User, err error) {
	return nil, nil
}

func (r Repository) Forget(ctx context.Context, token string) (err error) {
	return nil
}