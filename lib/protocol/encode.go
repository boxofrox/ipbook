package protocol

import (
	"encoding/json"
)

func encode(message interface{}) ([]byte, error) {
	return json.Marshal(message)
}

func encodeMessage (t Typeable) ([]byte, error) {
    data, err := encode(t)

    if nil != err {
        return nil, err
    }

    return encode(&Message{t.GetType(), data})
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
