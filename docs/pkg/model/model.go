package model

import "encoding/json"

type (
	Meta struct {
		ID        string            `json:"id"`
		Oid       uint32            `json:"-"`
		Name      string            `json:"name"`
		Mime      string            `json:"mime"`
		File      bool              `json:"file"`
		Public    bool              `json:"public"`
		Created   string            `json:"created"`
		Ts        int64             `json:"-"`
		Size      int               `json:"-"`
		Grant     []string          `json:"grant"`
	}

	SaveMeta struct {
		ID        string            `json:"id"`
		Name      string            `json:"name"`
		File      bool              `json:"file"`
		Public    bool              `json:"public"`
		Token     string            `json:"token"`
		Mime      string            `json:"mime"`
		Grant     []string          `json:"grant"`
	}

	ListMeta struct {
		Token     string            `json:"token"`
		Login     string            `json:"login"`
		Key       string            `json:"key"`
		Value     string            `json:"value"`
		Limit     int               `json:"limit"`
	}

	DeleteResponse map[string]interface{}

	ListResponse struct {
		Docs    []*Meta             `json:"docs"`
	}

	SaveResponse struct {
		JSON    json.RawMessage     `json:"json,omitempty"`
		File    string              `json:"file,omitempty"`
	}
)