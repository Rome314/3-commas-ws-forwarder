package api

import (
	"bytes"
	"encoding/json"
	"io"
)

type Message struct {
	Type       string      `json:"type"`
	Identifier interface{} `json:"identifier,omitempty"`
	Message    interface{} `json:"message,omitempty"`
}

func (m *Message) Reader() io.Reader {
	bts, _ := json.Marshal(m)
	return bytes.NewReader(bts)
}
