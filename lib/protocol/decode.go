package protocol

import (
	"encoding/json"
	"fmt"
)

func Decode(data []byte) (Messager, error) {
	var msg Message

	if err := json.Unmarshal(data, &msg); nil != err {
		return nil, err
	}

	switch msg.Type {
	case TYPE_ERROR_RESPONSE:
		return unmarshal(msg.Data, &ErrorResponse{})

	case TYPE_GET_IP_REQUEST:
		return unmarshal(msg.Data, &GetIpRequest{})

	case TYPE_GET_IP_RESPONSE:
		return unmarshal(msg.Data, &GetIpResponse{})

	case TYPE_SET_IP_REQUEST:
		return unmarshal(msg.Data, &SetIpRequest{})

	case TYPE_SET_IP_RESPONSE:
		return unmarshal(msg.Data, &SetIpResponse{})
	}

	return nil, &DecodingError{UNKNOWN_MESSAGE, msg.Data, nil}
}

func unmarshal(data []byte, v Messager) (Messager, error) {
	if err := json.Unmarshal(data, &v); nil != err {
		return nil, &DecodingError{DECODING_FAILED, data, err}
	}

	return v, nil
}

type DecodingError struct {
	Type     int
	Json     []byte
	Previous error
}

const (
	UNKNOWN_MESSAGE int = iota
	DECODING_FAILED
)

func (d DecodingError) Error() string {
	switch d.Type {
	case UNKNOWN_MESSAGE:
		return "cannot decode unknown message\n" +
			d.Previous.Error() +
			"\n---8<---8<---\n" +
			string(d.Json) +
			"\n---8<---8<---"

	case DECODING_FAILED:
		return fmt.Sprintf("decoding failed\n  %s\n  ---8<---8<---\n  %s\n  ---8<---8<---",
			d.Previous.Error(), string(d.Json))
	}

	return "unknown error"
}
