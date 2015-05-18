package protocol

import "fmt"

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

func (r *ErrorResponse) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *ErrorResponse) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("(%d) %s", r.Code, r.Reason)
}
