package protocol

import (
	"encoding/json"
)

type Typeable interface {
	GetType() int
}

type Message struct {
	Type int             `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (m *Message) DecodeData(dest interface{}) error {
	if err := json.Unmarshal(m.Data, &dest); nil != err {
		return &DecodingError{DECODING_FAILED, m.Data, err}
	}

	return nil
}

func (m *Message) ReadMessage(b []byte) error {
	if err := json.Unmarshal(b, m); nil != err {
		return &DecodingError{DECODING_FAILED, b, err}
	}

	return nil
}
