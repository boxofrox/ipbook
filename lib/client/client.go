package client

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	raddr   *net.UDPAddr
	timeout time.Duration
}

func New(address string, port int, timeout time.Duration) (*Client, error) {
	udp, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if nil != err {
		return nil, &Error{err}
	}

	return &Client{udp, timeout}, nil
}
