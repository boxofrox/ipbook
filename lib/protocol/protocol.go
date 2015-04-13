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

func SendErrorResponse(conn *net.UDPConn, addr *net.UDPAddr, code int, reason string) (int, error) {
	return SendMessage(conn, addr, &ErrorResponse{code, reason})
}

func SendGetIpRequest(conn *net.UDPConn, addr *net.UDPAddr, name string) (int, error) {
	return SendMessage(conn, addr, &GetIpRequest{name})
}

func SendGetIpResponse(conn *net.UDPConn, addr *net.UDPAddr, name, ip string) (int, error) {
	return SendMessage(conn, addr, &GetIpResponse{name, ip})
}

func SendMessage(conn *net.UDPConn, addr *net.UDPAddr, msg messager) (int, error) {
	if nil == conn {
		return 0, errors.New("socket is nil")
	}

	if nil == addr {
		return 0, errors.New("address is nil")
	}

	payload, err := encode(msg)

	if nil != err {
		return 0, err
	}

	return conn.WriteToUDP(payload, addr)
}

func SendSetIpRequest(conn *net.UDPConn, addr *net.UDPAddr, name, ip string) (int, error) {
	return SendMessage(conn, addr, &SetIpRequest{name, ip})
}

func SendSetIpResponse(conn *net.UDPConn, addr *net.UDPAddr, status Status, msg string) (int, error) {
	return SendMessage(conn, addr, &SetIpResponse{status, msg})
}
