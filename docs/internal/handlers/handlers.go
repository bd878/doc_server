package handlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/rs/zerolog"
	server "github.com/bd878/doc_server/pkg/model"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type UsersGateway interface {
	Auth(ctx context.Context, token string) (login string, err error)
}

type Controller interface {
	List(ctx context.Context, owner, login, key, value string, limit int) (docs []*docs.Meta, err error)
	Save(ctx context.Context, owner string, f multipart.File, json []byte, meta *docs.Meta) (err error)
	GetMeta(ctx context.Context, id, login string) (doc *docs.Meta, err error)
	ReadJSON(ctx context.Context, id string) (json json.RawMessage, err error)
	ReadFileStream(ctx context.Context, oid uint32, w io.Writer) (err error)
	Delete(ctx context.Context, id string) (err error)
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
	mux.HandleFunc("GET     /api/docs/{id}", h.Get)
	mux.HandleFunc("HEAD    /api/docs/{id}", h.GetHead)
	mux.HandleFunc("DELETE  /api/docs/{id}", h.Delete)
}

func (h handlers) Save(w http.ResponseWriter, req *http.Request) {
	var meta docs.SaveMeta

	err := req.ParseMultipartForm(10 << 20 /* 10 MB */)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse form")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoForm,
				Text: "failed to parse form",
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
		h.logger.Error().Err(err).Msg("failed to unmarshal meta")
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

	login, err := h.gateway.Auth(req.Context(), meta.Token)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Text: "not authorized",
			},
		})
		return
	}

	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var f multipart.File
	if meta.File {
		f, _, err = req.FormFile("file")
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to read form file")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Error: &server.ErrorCode{
					Code: docs.CodeNoFile,
					Text: "file required",
				},
			})
			return
		}
	}

	jsonData := req.PostFormValue("json")
	if !meta.File && jsonData == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeNoJSON,
				Text: "json required",
			},
		})
		return
	}

	err = h.ctrl.Save(req.Context(), login, f, []byte(jsonData), &docs.Meta{
		Name:     meta.Name,
		File:     meta.File,
		Mime:     meta.Mime,
		Public:   meta.Public,
		Grant:    meta.Grant,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to save file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(docs.SaveResponse{
		File: meta.Name,
		JSON: json.RawMessage(jsonData),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal response")
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
		h.logger.Error().Err(err).Msg("failed to parse form")
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

	owner, err := h.gateway.Auth(req.Context(), token)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Text: "not authorized",
			},
		})
		return
	}

	if owner == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	login := req.FormValue("login")

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
		h.logger.Error().Err(err).Msg("failed to parse limit")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeBadLimit,
				Text: "bad limit param",
			},
		})
		return
	}

	list, err := h.ctrl.List(req.Context(), owner, login, key, value, limit)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(docs.ListResponse{
		Docs: list,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal response")
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
	err := req.ParseForm()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse form")
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

	login, err := h.gateway.Auth(req.Context(), token)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Text: "not authorized",
			},
		})
		return
	}

	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := req.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	meta, err := h.ctrl.GetMeta(req.Context(), id, login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeDocNotFound,
				Text: "document not found",
			},
		})
		return
	}

	if meta.File {
		w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + meta.Name)
		w.Header().Set("Content-Type", meta.Mime)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", meta.Size))
		w.Header().Set("Date", meta.Created)

		err = h.ctrl.ReadFileStream(req.Context(), meta.Oid, w)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to read file stream")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", meta.Size))
		w.Header().Set("Date", meta.Created)

		jsonData, err := h.ctrl.ReadJSON(req.Context(), id)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to read json data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(server.ServerResponse{
			Data: jsonData,
		})
		return
	}
}

func (h handlers) GetHead(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse form")
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

	login, err := h.gateway.Auth(req.Context(), token)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to auth")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := req.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	meta, err := h.ctrl.GetMeta(req.Context(), id, login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: docs.CodeDocNotFound,
				Text: "document not found",
			},
		})
		return
	}

	if meta.File {
		w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + meta.Name)
		w.Header().Set("Content-Type", meta.Mime)
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", meta.Size))
	w.Header().Set("Date", meta.Created)
}

func (h handlers) Delete(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse form")
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

	_, err = h.gateway.Auth(req.Context(), token)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Text: "not authorized",
			},
		})
		return
	}

	id := req.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ctrl.Delete(req.Context(), id)
	if err != nil {
		switch err {
		case docs.ErrNoDoc:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Error: &server.ErrorCode{
					Code: docs.CodeDocNotFound,
					Text: "no document",
				},
			})
			return
		default:
			h.logger.Error().Err(err).Msg("failed to delete")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	response := json.RawMessage([]byte(fmt.Sprintf(`{"%s": true}`, id)))
	json.NewEncoder(w).Encode(server.ServerResponse{
		Response: response,
	})
}
