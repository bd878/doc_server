package handlers

import (
	"context"
	"net/http"
	"github.com/bd878/doc_server/users/pkg/model"
)

type Controller interface {
	Register(ctx context.Context, login, password string) (err error)
	Auth(ctx context.Context, login, password string) (user *model.User, err error)
	Logout(ctx context.Context, token string) (err error)
}

type handlers struct {
	ctrl Controller
}

func RegisterHandlers(mux *http.ServeMux, ctrl Controller) {
	h := &handlers{ctrl}

	mux.HandleFunc("/api/register", h.Register)
	mux.HandleFunc("/api/auth", h.Auth)
	mux.HandleFunc("/api/auth/:token", h.Logout)
}

func (h handlers) Register(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handlers) Auth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handlers) Logout(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
