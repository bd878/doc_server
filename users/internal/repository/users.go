package repository

import (
	"fmt"
	"errors"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/doc_server/users/pkg/model"
)

type Repository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func New(tableName string, pool *pgxpool.Pool) *Repository {
	return &Repository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r Repository) Save(ctx context.Context, token, login, hashedPassword string) (err error) {
	const query = "INSERT INTO %s(login, salt, token) VALUES ($1, $2, $3)"

	_, err = r.pool.Exec(ctx, r.table(query), login, hashedPassword, token)

	return
}

func (r Repository) Find(ctx context.Context, login, token string) (user *model.User, err error) {
	const query = "SELECT token, login, salt FROM %s WHERE (token = $1 OR login = $2)"

	user = &model.User{
	}

	var nullToken *string

	err = r.pool.QueryRow(ctx, r.table(query), token, login).Scan(&nullToken, &user.Login, &user.HashedPassword)
	if err != nil {
		return nil, err
	}

	if nullToken != nil {
		user.Token = *nullToken
	}

	return
}

func (r Repository) Forget(ctx context.Context, token string) (err error) {
	const query = "UPDATE %s SET token = NULL WHERE token = $1"

	result, err := r.pool.Exec(ctx, r.table(query), token)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows != 1 {
		return errors.New("no user")
	}

	return
}

func (r Repository) Auth(ctx context.Context, login, token string) (err error) {
	const query = "UPDATE %s SET token = $2 WHERE login = $1"

	result, err := r.pool.Exec(ctx, r.table(query), login, token)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows != 1 {
		return errors.New("no user")
	}

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}