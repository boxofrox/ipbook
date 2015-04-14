package protocol

import (
	"errors"
	"net"

	boxnet "github.com/boxofrox/ipbook/lib/net"
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

func SendErrorResponse(conn *boxnet.Conn, addr *net.UDPAddr, code int, reason string) error {
	return SendMessage(conn, addr, &ErrorResponse{code, reason})
}

func SendGetIpRequest(conn *boxnet.Conn, addr *net.UDPAddr, name string) error {
	return SendMessage(conn, addr, &GetIpRequest{name})
}

func SendGetIpResponse(conn *boxnet.Conn, addr *net.UDPAddr, name, ip string) error {
	return SendMessage(conn, addr, &GetIpResponse{name, ip})
}

func SendMessage(conn *boxnet.Conn, addr *net.UDPAddr, msg Messager) error {
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

	_, err = conn.WriteToUDP(payload, addr)
	return err
}

func SendSetIpRequest(conn *boxnet.Conn, addr *net.UDPAddr, name, ip string) error {
	return SendMessage(conn, addr, &SetIpRequest{name, ip})
}

func SendSetIpResponse(conn *boxnet.Conn, addr *net.UDPAddr, status Status, msg string) error {
	return SendMessage(conn, addr, &SetIpResponse{status, msg})
}
