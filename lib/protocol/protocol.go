package protocol

import (
	"errors"
	"net"
)

const (
	TYPE_ERROR_RESPONSE = iota
	TYPE_GET_IP_REQUEST
	TYPE_GET_IP_RESPONSE
	TYPE_SET_IP_REQUEST
	TYPE_SET_IP_RESPONSE
)

func ReadGetIpResponse(buffer []byte) (*GetIpResponse, error) {
	msg, err := Decode(buffer)

	if nil != err {
		return nil, err
	}

	return msg.(*GetIpResponse), nil
}

func ReadSetIpResponse(buffer []byte) (*SetIpResponse, error) {
	msg, err := Decode(buffer)

	if nil != err {
		return nil, err
	}

	return msg.(*SetIpResponse), nil
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

func SendMessage(conn net.PacketConn, addr net.Addr, msg Messager) error {
	if nil == conn {
		return errors.New("socket is nil")
	}

	if nil == addr {
		return errors.New("address is nil")
	}

	payload, err := encode(msg)

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
