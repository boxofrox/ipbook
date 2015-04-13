package protocol

import (
	"encoding/json"
)

type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

const (
	BAD_REQUEST = iota
	NAME_NOT_FOUND
	INVALID_NAME
)

func (resp *ErrorResponse) GetType() int {
	return TYPE_ERROR_RESPONSE
}

func (resp *ErrorResponse) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*resp)
	if nil != err {
		return nil, err
	}

	msg := Message{
		Type: TYPE_ERROR_RESPONSE,
		Data: data,
	}
	return json.Marshal(&msg)
}
