package protocol

import (
	"encoding/json"
)

type Messager interface {
	GetType() int
}

type Message struct {
	Type int             `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (m *Message) GetType() int {
	return m.Type
}
