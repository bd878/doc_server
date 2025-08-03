package repository

import (
	"io"
	"fmt"
	"context"
	"encoding/json"
	"mime/multipart"
	"github.com/rs/zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Repository struct {
	tableName              string
	permissionsTableName   string
	pool                  *pgxpool.Pool
	log                    zerolog.Logger
}

func New(log zerolog.Logger, tableName, permissionsTableName string, pool *pgxpool.Pool) *Repository {
	return &Repository{
		log:                   log,
		tableName:             tableName,
		permissionsTableName:  permissionsTableName,
		pool:                  pool,
	}
}

func (r *Repository) SaveFile(ctx context.Context, owner string, f multipart.File, meta *docs.Meta) (err error) {
	const query = "INSERT INTO %s(id, oid, name, file, public, mime, owner_login) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	const queryPermissions = "INSERT INTO %s(file_id, user_login) VALUES ($1, $2)"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	lb := tx.LargeObjects()
	oid, err := lb.Create(ctx, 0)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to create large object")
		return err
	}

	object, err := lb.Open(ctx, oid, pgx.LargeObjectModeRead)
	if err != nil {
		return err
	}

	_, err = io.Copy(object, f)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(query), meta.ID, oid, meta.Name, meta.File, meta.Public, meta.Mime, owner)
	if err != nil {
		return
	}

	for _, login := range meta.Grant {
		_, err = tx.Exec(ctx, r.permissionsTable(queryPermissions), meta.ID, login)
		if err != nil {
			return
		}
	}

	err = tx.Commit(ctx)

	return
}

func (r *Repository) SaveJSON(ctx context.Context, owner string, data []byte, meta *docs.Meta) (err error) {
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

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func (r Repository) permissionsTable(query string) string {
	return fmt.Sprintf(query, r.permissionsTableName)
}