package handlers

import (
	"context"
	"mime/multipart"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/rs/zerolog"
	server "github.com/bd878/doc_server/pkg/model"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type UsersGateway interface {
	Auth(ctx context.Context, token string) (ok bool, err error)
}

type Controller interface {
	Save(ctx context.Context, f multipart.File, meta *docs.Meta) (err error)
	List(ctx context.Context, key, value string, limit int) (docs []*docs.Meta, isLastPage bool, err error)
	Get(ctx context.Context, id int) (doc *docs.Meta, err error)
	Delete(ctx context.Context, id int) (err error)
}

type handlers struct {
	ctrl    Controller
	logger  zerolog.Logger
	gateway UsersGateway
}

func RegisterHandlers(mux *http.ServeMux, ctrl Controller, gateway UsersGateway, logger zerolog.Logger) {
	h := &handlers{ctrl, logger, gateway}

	mux.HandleFunc("POST    /api/docs", h.Save)
	mux.HandleFunc("GET     /api/docs", h.List)
	mux.HandleFunc("HEAD    /api/docs", h.ListHead)
	mux.HandleFunc("GET     /api/docs/:id", h.Get)
	mux.HandleFunc("HEAD    /api/docs/:id", h.GetHead)
	mux.HandleFunc("DELETE  /api/docs/:id", h.Delete)
}

func (h handlers) Save(w http.ResponseWriter, req *http.Request) {
	var meta docs.SaveMeta

	err := req.ParseMultipartForm(10 << 20 /* 10 MB */)
	if err != nil {
		h.logger.Error().Err(err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeRequestTooLarge,
				Text: "request too large",
			},
		})
		return
	}

	rawMeta := req.PostFormValue("meta")
	if rawMeta == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeNoMeta,
				Text: "meta required",
			},
		})
		return
	}

	err = json.Unmarshal([]byte(rawMeta), &meta)
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if meta.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoToken,
				Text: "no token",
			},
		})
		return
	}

	ok, err := h.gateway.Auth(req.Context(), meta.Token)
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonData := req.PostFormValue("jsonData")

	f, _, err := req.FormFile("file")
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeNoFile,
				Text: "file required",
			},
		})
		return
	}

	err = h.ctrl.Save(req.Context(), f, &docs.Meta{
		Name:     meta.Name,
		File:     meta.File,
		Mime:     meta.Mime,
		Public:   meta.Public,
		Grant:    meta.Grant,
	})
	if err != nil {
		h.logger.Error().Err(err)
		switch err {
		case server.ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	response, err := json.Marshal(docs.SaveResponse{
		JSON: json.RawMessage([]byte(jsonData)),
		File: meta.Name,
	})
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Data: json.RawMessage(response),
	})
}

func (h handlers) List(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := req.FormValue("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoToken,
				Text: "no token",
			},
		})
		return
	}

	ok, err := h.gateway.Auth(req.Context(), token)
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeNoKey,
				Text: "no key",
			},
		})
		return
	}

	value := req.FormValue("value")
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeNoValue,
				Text: "no value",
			},
		})
		return
	}

	rawLimit := req.FormValue("limit")
	if rawLimit == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeNoLimit,
				Text: "no limit",
			},
		})
		return
	}

	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeBadLimit,
				Text: "bad limit param",
			},
		})
		return
	}

	list, _, err := h.ctrl.List(req.Context(), key, value, limit)
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(docs.ListResponse{
		Docs: list,
	})
	if err != nil {
		h.logger.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Data: json.RawMessage(response),
	})
}

func (h handlers) ListHead(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) Get(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) GetHead(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) Delete(w http.ResponseWriter, req *http.Request) {
}
