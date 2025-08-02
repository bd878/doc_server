package controller

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/bd878/doc_server/users/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, token, login, hashedPassword string) (err error)
	Find(ctx context.Context, login, token string) (user *model.User, err error)
	Forget(ctx context.Context, token string) (err error)
}

type Controller struct {
	repo  Repository
	token string
}

func New(repo Repository, token string) *Controller {
	return &Controller{repo, token}
}

func (r Controller) Register(ctx context.Context, adminToken, login, password string) (err error) {
	if adminToken != r.token {
		return model.ErrWrongToken
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	token := uuid.New().String()

	err = r.repo.Save(ctx, token, login, string(hashed))

	return
}

func (r Controller) Auth(ctx context.Context, token string) (user *model.User, err error) {
	return r.repo.Find(ctx, "", token)
}

func (r Controller) Login(ctx context.Context, login, password string) (user *model.User, err error) {
	user, err = r.repo.Find(ctx, login, "")
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	return
}

func (r Controller) Logout(ctx context.Context, token string) (err error) {
	return r.repo.Forget(ctx, token)
}