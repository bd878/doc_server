package controller

import (
	"context"
	"github.com/bd878/doc_server/users/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, user *model.User) (err error)
	Find(ctx context.Context, token string) (user *model.User, err error)
	Forget(ctx context.Context, token string) (err error)
}

type Controller struct {
	repo Repository
}

func New(repo Repository) *Controller {
	return &Controller{repo}
}

func (r Controller) Register(ctx context.Context, login, password string) (err error) {
	return nil
}

func (r Controller) Auth(ctx context.Context, login, password string) (user *model.User, err error) {
	return nil, nil
}

func (r Controller) Logout(ctx context.Context, token string) (err error) {
	return nil
}