package repository

import (
	"io"
	"os"
	"fmt"
	"time"
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
	pool                  *pgxpool.Pool
	log                    zerolog.Logger
}

func New(log zerolog.Logger, tableName string, pool *pgxpool.Pool) *Repository {
	return &Repository{
		log:                   log,
		tableName:             tableName,
		pool:                  pool,
	}
}

func (r *Repository) Save(ctx context.Context, owner string, f multipart.File, jsonData []byte, meta *docs.Meta) (err error) {
	const query = "INSERT INTO %s(id, oid, name, file, json, public, mime, owner_login, grant_logins) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

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
			fmt.Fprintf(os.Stderr, "rollback with error: %w", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var oid uint32
	if meta.File && f != nil {
		lb := tx.LargeObjects()
		oid, err = lb.Create(ctx, 0)
		if err != nil {
			r.log.Error().Err(err).Msg("failed to create large object")
			return err
		}

		object, err := lb.Open(ctx, oid, pgx.LargeObjectModeWrite)
		if err != nil {
			return err
		}
		defer object.Close()

		_, err = io.Copy(object, f)
		if err != nil {
			return err
		}
	}

	grant, err := json.Marshal(meta.Grant)
	if err != nil {
		return
	}

	// no file, only json
	if oid == 0 {
		_, err = tx.Exec(ctx, r.table(query), meta.ID, nil, meta.Name, false, jsonData, meta.Public, meta.Mime, owner, grant)
		if err != nil {
			return
		}
	} else {
		_, err = tx.Exec(ctx, r.table(query), meta.ID, oid, meta.Name, true, jsonData, meta.Public, meta.Mime, owner, grant)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) List(ctx context.Context, owner, login, key, value string, limit int) (list []*docs.Meta, err error) {
	const queryLogin = "SELECT id, name, file, public, mime, created_at, grant_logins FROM %s, jsonb_array_elements_text(grant_logins) AS login WHERE %s = $2 AND login = $1 LIMIT $3"
	const query = "SELECT id, name, file, public, mime, created_at, grant_logins FROM %s WHERE %s = $2 AND owner_login = $1 LIMIT $3"

	r.log.Log().Str("owner", owner).Str("login", login).Str("key", key).Str("value", value).Int("limit", limit).Msg("list docs")

	var rows pgx.Rows
	if login == "" {
		rows, err = r.pool.Query(ctx, fmt.Sprintf(query, r.tableName, key), owner, value, limit)
		if err != nil {
			return
		}
	} else {
		rows, err = r.pool.Query(ctx, fmt.Sprintf(queryLogin, r.tableName, key), login, value, limit)
		if err != nil {
			return
		}
	}
	defer rows.Close()

	list = make([]*docs.Meta, 0)
	for rows.Next() {
		meta := &docs.Meta{
		}

		var grant []byte
		var created time.Time

		err = rows.Scan(&meta.ID, &meta.Name, &meta.File, &meta.Public, &meta.Mime, &created, &grant)
		if err != nil {
			return
		}

		if grant != nil {
			err = json.Unmarshal(grant, &meta.Grant)
			if err != nil {
				return
			}
		}

		meta.Created = created.Format(time.DateTime)

		list = append(list, meta)
	}

	if err = rows.Err(); err != nil {
		return
	}

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
