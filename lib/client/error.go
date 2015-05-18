package client

import (
	"net"

	"github.com/boxofrox/ipbook/lib/protocol"
)

type Error struct {
	Err error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) ErrorResponse() bool {
	switch e.Err.(type) {
	case *protocol.ErrorResponse:
		return true
	default:
		return false
	}
}

func (e *Error) Temporary() bool {
	switch e.Err.(type) {
	case *net.OpError:
		return e.Err.(*net.OpError).Temporary()
	default:
		return false
	}
}

func (e *Error) Timeout() bool {
	switch e.Err.(type) {
	case *net.OpError:
		return e.Err.(*net.OpError).Timeout()
	default:
		return false
	}
}
