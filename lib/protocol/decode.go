package protocol

import (
	"encoding/json"
	"fmt"
)

func decode (b []byte, v interface{}) error {
	if err := json.Unmarshal(b, v); nil != err {
		return &DecodingError{DECODING_FAILED, b, err}
	}

	return nil
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
