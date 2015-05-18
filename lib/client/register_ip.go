package client

import (
	"fmt"
	"net"
	"time"

	"github.com/boxofrox/ipbook/lib/protocol"
)

func (c *Client) RegisterIp(name string, ip string) error {
	var (
		conn net.PacketConn
		err  error
		p    *protocol.Packet
		m    *protocol.Message
	)

	if conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}); nil != err {
		return &Error{err}
	}

	if 0 < c.timeout {
		conn.SetDeadline(time.Now().Add(c.timeout))
	}

	if err = protocol.SendSetIpRequest(conn, c.raddr, name, ip); nil != err {
		return &Error{err}
	}

	if p, err = protocol.ReadPacket(conn); nil != err {
		return &Error{err}
	}

	if m, err = p.ReadMessage(); nil != err {
		return &Error{err}
	}

	switch m.Type {
	case protocol.TYPE_ERROR_RESPONSE:
		r := &protocol.ErrorResponse{}

		if err = r.ReadFrom(m); nil != err {
			return &Error{err}
		}

		return &Error{r}

	case protocol.TYPE_SET_IP_RESPONSE:
		r := &protocol.SetIpResponse{}

		if err = r.ReadFrom(m); nil != err {
			return &Error{err}
		}

		return nil

	default:
		return &Error{fmt.Errorf("unknown response type: %d", m.Type)}
	}
}
