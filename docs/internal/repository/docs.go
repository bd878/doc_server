package repository

import (
	"io"
	"os"
	"fmt"
	"time"
	"context"
	"errors"
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
	const queryLogin = "SELECT id, oid, name, file, public, mime, created_at, grant_logins FROM %s, jsonb_array_elements_text(grant_logins) AS login WHERE %s = $2 AND login = $1 ORDER BY created_at DESC LIMIT $3"
	const query = "SELECT id, oid, name, file, public, mime, created_at, grant_logins FROM %s WHERE %s = $2 AND owner_login = $1 ORDER BY created_at DESC LIMIT $3"

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
		var oid *uint32

		err = rows.Scan(&meta.ID, &oid, &meta.Name, &meta.File, &meta.Public, &meta.Mime, &created, &grant)
		if err != nil {
			return
		}

		if grant != nil {
			err = json.Unmarshal(grant, &meta.Grant)
			if err != nil {
				return
			}
		}

		if oid != nil {
			meta.Oid = *oid
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
	const query = "SELECT id, oid, name, file, public, mime, created_at, grant_logins FROM %s WHERE id = $1"

	r.log.Log().Str("id", id).Msg("get meta")

	var grant []byte
	var created time.Time
	var oid *uint32

	meta = &docs.Meta{}

	err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&meta.ID, &oid, &meta.Name, &meta.File, &meta.Public, &meta.Mime, &created, &grant)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			err = docs.ErrNoDoc
		}
		return nil, err
	}

	if grant != nil {
		err = json.Unmarshal(grant, &meta.Grant)
		if err != nil {
			return
		}
	}

	if oid != nil {
		meta.Oid = *oid
	}

	meta.Created = created.Format(time.DateTime)

	return
}

func (r *Repository) ReadFile(ctx context.Context, oid uint32, writer io.Writer) (err error) {
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

	lb := tx.LargeObjects()
	object, err := lb.Open(ctx, oid, pgx.LargeObjectModeRead)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}

	object.Close()

	return
}

func (r *Repository) ReadJSON(ctx context.Context, id string) (result json.RawMessage, err error) {
	const query = "SELECT json FROM %s WHERE id = $1"

	var jsonData []byte

	err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&jsonData)
	if err != nil {
		return
	}

	return json.RawMessage(jsonData), nil
}

func (r *Repository) Delete(ctx context.Context, id string) (err error) {
	const query = "SELECT oid FROM %s WHERE id = $1"
	const deleteQuery = "DELETE FROM %s WHERE id = $1"

	r.log.Log().Str("id", id).Msg("delete file")

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
			if errors.Is(err, docs.ErrNoDoc) {
				return
			} else {
				fmt.Fprintf(os.Stderr, "rollback with error: %w", err)
				err = tx.Rollback(ctx)
			}
		default:
			err = tx.Commit(ctx)
		}
	}()

	var oid int

	err = tx.QueryRow(ctx, r.table(query), id).Scan(&oid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return docs.ErrNoDoc
		}
		return
	}

	lb := tx.LargeObjects()
	err = lb.Unlink(ctx, uint32(oid))
	if err != nil {
		return
	}

	result, err := tx.Exec(ctx, r.table(deleteQuery), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = docs.ErrNoDoc
		}
		return err
	}

	if result.RowsAffected() != 1 {
		return docs.ErrNoDoc
	}

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
