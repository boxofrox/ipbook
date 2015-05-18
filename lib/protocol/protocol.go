package protocol

import (
	"fmt"
	"net"

	"github.com/boxofrox/ipbook/lib/buffer"
	"github.com/boxofrox/ipbook/lib/pool"
)

const (
	TYPE_ERROR_RESPONSE = iota
	TYPE_GET_IP_REQUEST
	TYPE_GET_IP_RESPONSE
	TYPE_SET_IP_REQUEST
	TYPE_SET_IP_RESPONSE
	TYPE_SET_PUBLIC_IP_REQUEST
	TYPE_LAST // not a valid message type
)

// global buffer pool for protocol library
var bufferPool = pool.New(5, buffer.CreateUdpBuffer)

func IsValidMessage(msgType int) bool {
	return 0 <= msgType && msgType < TYPE_LAST
}

func SendErrorResponse(conn net.PacketConn, addr net.Addr, code int, reason string) error {
	return SendMessage(conn, addr, &ErrorResponse{code, reason})
}

func SendGetIpRequest(conn net.PacketConn, addr net.Addr, name string) error {
	return SendMessage(conn, addr, &GetIpRequest{name})
}

func SendGetIpResponse(conn net.PacketConn, addr net.Addr, name, ip string) error {
	return SendMessage(conn, addr, &GetIpResponse{name, ip})
}

func SendMessage(conn net.PacketConn, addr net.Addr, data Encodable) error {
	if nil == conn {
		return fmt.Errorf("socket is nil")
	}

	if nil == addr {
		return fmt.Errorf("address is nil")
	}

	payload, err := data.EncodeMessage()

	if nil != err {
		return err
	}

	_, err = conn.WriteTo(payload, addr)
	return err
}

func SendSetIpRequest(conn net.PacketConn, addr net.Addr, name, ip string) error {
	return SendMessage(conn, addr, &SetIpRequest{name, ip})
}

func SendSetIpResponse(conn net.PacketConn, addr net.Addr, status Status, msg string) error {
	return SendMessage(conn, addr, &SetIpResponse{status, msg})
}

func SendSetPublicIpRequest(conn net.PacketConn, addr net.Addr, name string) error {
	return SendMessage(conn, addr, &SetPublicIpRequest{name})
}
