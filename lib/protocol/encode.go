package protocol

import (
	"encoding/json"
)

func encode(message messager) ([]byte, error) {
	return json.Marshal(message)
}

type EncodingError int

const (
	NOT_A_MESSAGE_TYPE EncodingError = iota
)

func (e EncodingError) Error() string {
	switch e {
	case NOT_A_MESSAGE_TYPE:
		return "argument is not a valid message type"
	}

	return "unknown error"
}

