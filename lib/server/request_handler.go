package server

import (
	"net"

	"github.com/boxofrox/ipbook/lib/protocol"
)

type RequestHandler func(s *Server, addr net.Addr, msg *protocol.Message)
