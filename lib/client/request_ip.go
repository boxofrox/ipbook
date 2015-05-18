package client

import (
	"fmt"
	"net"
	"time"

	"github.com/boxofrox/ipbook/lib/protocol"
)

func (c *Client) RequestIp(name string) (string, error) {
	var (
		conn net.PacketConn
		p    *protocol.Packet
		m    *protocol.Message
		err  error
	)

	if conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}); nil != err {
		return "", &Error{err}
	}

	if 0 < c.timeout {
		conn.SetDeadline(time.Now().Add(c.timeout))
	}

	if err = protocol.SendGetIpRequest(conn, c.raddr, name); nil != err {
		return "", &Error{err}
	}

	if p, err = protocol.ReadPacket(conn); nil != err {
		return "", &Error{err}
	}

	if m, err = p.ReadMessage(); nil != err {
		return "", &Error{err}
	}

	switch m.Type {
	case protocol.TYPE_ERROR_RESPONSE:
		r := &protocol.ErrorResponse{}

		if err = r.ReadFrom(m); nil != err {
			return "", &Error{err}
		}

		return "", &Error{r}

	case protocol.TYPE_GET_IP_RESPONSE:
		r := &protocol.GetIpResponse{}

		if err = r.ReadFrom(m); nil != err {
			return "", &Error{err}
		}

		return r.Ip, nil

	default:
		err = fmt.Errorf("expecting GetIpResponse (%d), got (%d)",
			protocol.TYPE_GET_IP_RESPONSE,
			m.Type)
		return "", &Error{err}
	}
}
