package model

import "encoding/json"

type (
	Doc struct {
		Meta   Meta
		File   File
	}

	Meta struct {
		ID        string            `json:"id"`
		Name      string            `json:"name"`
		File      bool              `json:"file"`
		Public    bool              `json:"public"`
		Token    *string            `json:"token,omitempty"`
		Mime      string            `json:"mime"`
		Grant     []string          `json:"grant"`
	}

	File struct {
		JSON   json.RawMessage   `json:"json"`
		File   []byte            `json:"file"`
	}

	DeleteResponse map[string]interface{}

	SaveResponse struct {
		JSON    json.RawMessage   `json:"json,omitempty"`
		File    string            `json:"file"`
	}
)