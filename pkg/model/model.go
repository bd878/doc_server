package model

import "encoding/json"

type (
	ErrorCode struct {
		Code    int         `json:"code"`
		Text    string      `json:"text"`
	}

	ServerResponse struct {
		Error     *ErrorCode         `json:"error,omitempty"`
		Response   json.RawMessage   `json:"response,omitempty"`
		Data       json.RawMessage   `json:"data,omitempty"`
	}
)
