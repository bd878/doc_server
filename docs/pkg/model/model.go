package model

import "encoding/json"

type (
	Doc struct {
		ID     string            `json:"id"`
		Name   string            `json:"name"`
		File   bool              `json:"file"`
		Public bool              `json:"public"`
		Token  string            `json:"token"`
		Mime   string            `json:"mime"`
		Grant  []string          `json:"grant"`
	}

	NewDocData struct {
		JSON   json.RawMessage   `json:"json"`
		File   string            `json:"file"`
	}

	DeleteResponse map[string]interface{}
)