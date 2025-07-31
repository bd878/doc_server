package handlers

import (
	"io"
	"context"
	"net/http"
	"encoding/json"
	"github.com/bd878/doc_server/users/pkg/model"
)

type Controller interface {
	Register(ctx context.Context, adminToken, login, password string) (err error)
	Auth(ctx context.Context, login, password string) (user *model.User, err error)
	Logout(ctx context.Context, token string) (err error)
}

type handlers struct {
	ctrl Controller
}

func RegisterHandlers(mux *http.ServeMux, ctrl Controller) {
	h := &handlers{ctrl}

	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/auth", h.Auth)
	mux.HandleFunc("DELETE /api/auth/:token", h.Logout)
}

func verifyPassword(password string) (eightOrMore, upper, number, symbol bool) {
	return
}

func (h handlers) Register(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var body pkg.Register
	if err := json.Unmarshal(data, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eightOrMore, upper, number, symbol := verifyPassword(body.Password)
	if !eightOrMore {

	}
	if !upper {

	}
	if !number {

	}
	if !symbol {
		
	}

	err := h.ctrl.Register(req.Context(), body.Token, body.Login, body.Password)
	if err != nil {

	}
}

func (h handlers) Auth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handlers) Logout(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
