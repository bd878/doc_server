package handlers

import (
	"context"
	"net/http"
	"github.com/rs/zerolog"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Controller interface {
	Save(ctx context.Context, doc *docs.Doc) (err error)
	List(ctx context.Context, key, value string, limit int) (docs []*docs.Doc, isLastPage bool, err error)
	Get(ctx context.Context, id int) (doc *docs.Doc, err error)
	Delete(ctx context.Context, id int) (err error)
}

type handlers struct {
	ctrl   Controller
	logger zerolog.Logger
}

func RegisterHandlers(mux *http.ServeMux, ctrl Controller, logger zerolog.Logger) {
	h := &handlers{ctrl, logger}

	mux.HandleFunc("POST    /api/docs", h.Save)
	mux.HandleFunc("GET     /api/docs", h.List)
	mux.HandleFunc("HEAD    /api/docs", h.ListHead)
	mux.HandleFunc("GET     /api/docs/:id", h.Get)
	mux.HandleFunc("HEAD    /api/docs/:id", h.GetHead)
	mux.HandleFunc("DELETE  /api/docs/:id", h.Delete)
}

func (h handlers) Save(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) List(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) ListHead(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) Get(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) GetHead(w http.ResponseWriter, req *http.Request) {
}

func (h handlers) Delete(w http.ResponseWriter, req *http.Request) {
}
