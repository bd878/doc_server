package handlers

import (
	"io"
	"errors"
	"context"
	"unicode"
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog"
	server "github.com/bd878/doc_server/pkg/model"
	users "github.com/bd878/doc_server/users/pkg/model"
)

type Controller interface {
	Register(ctx context.Context, adminToken, login, password string) (err error)
	Auth(ctx context.Context, login, password string) (user *users.User, err error)
	Logout(ctx context.Context, token string) (err error)
}

type handlers struct {
	ctrl   Controller
	logger zerolog.Logger
}

func RegisterHandlers(mux *http.ServeMux, ctrl Controller, logger zerolog.Logger) {
	h := &handlers{ctrl, logger}

	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/auth", h.Auth)
	mux.HandleFunc("DELETE /api/auth/:token", h.Logout)
}

func verifyPassword(password string) (eightOrMore, twoLetters, oneNumber, oneSpecial bool) {
	if len(password) >= 8 {
		eightOrMore = true
	}

	var oneLower, oneUpper bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			oneNumber = true
		case unicode.IsUpper(c):
			oneUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			oneSpecial = true
		case unicode.IsLetter(c):
			oneLower = true
		default:
		}
	}
	twoLetters = oneUpper && oneLower
	return
}

func (h handlers) Register(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var body users.Register
	if err := json.Unmarshal(data, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eightOrMore, twoLetters, oneNumber, oneSpecial := verifyPassword(body.Password)
	if !eightOrMore {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: users.CodePasswordTooShort,
				Text: "password is less than 8 symbols",
			},
		})
		return
	}
	if !twoLetters {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: users.CodePasswordUpperLower,
				Text: "password must have upper und lower letter",
			},
		})
		return
	}
	if !oneNumber {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: users.CodePasswordOneNumber,
				Text: "password must have at least one number",
			},
		})
		return
	}
	if !oneSpecial {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: users.CodePasswordOneSpecial,
				Text: "password must have at least one special symbol",
			},
		})
		return
	}

	err = h.ctrl.Register(req.Context(), body.Token, body.Login, body.Password)
	if err != nil {
		h.logger.Error().Err(err)

		if errors.Is(err, users.ErrWrongToken) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Error: &server.ErrorCode{
					Code: users.CodeWrongToken,
					Text: "wrong admin token",
				},
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: users.CodeRegisterFailed,
				Text: "failed to register user",
			},
		})
		return
	}

	response, err := json.Marshal(users.RegisterResponse{
		Login: body.Login,
	})
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Response: json.RawMessage(response),
	})
}

func (h handlers) Auth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handlers) Logout(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
