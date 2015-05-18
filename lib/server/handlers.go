package server

import (
	"log"
	"net"

	"github.com/boxofrox/ipbook/lib/protocol"
)

func handleGetIpRequest(s *Server, addr net.Addr, m *protocol.Message) {
	var (
		err error
		ip  = onlyIp(addr)
	)

	r := &protocol.GetIpRequest{}

	if err = r.ReadFrom(m); nil != err {
		log.Printf("error: from host (%s): %s", ip, err)
		s.sendErrorResponse(addr, protocol.BAD_REQUEST, err.Error())
		return
	}

	if false == s.registry.Contains(r.Name) {

		log.Printf("Host (%s): requested name (%s) not found in registry.", ip, r.Name)
		s.sendErrorResponse(addr, protocol.NAME_NOT_FOUND, "name not found in registry")
		return
	}

	ip, _ := s.registry.Get(r.Name)

	if s.sendGetIpResponse(addr, r.Name, ip) {
		log.Printf("Host (%s): requested IP of (%s).", ip, r.Name)
	}
}

func handleSetIpRequest(s *Server, addr net.Addr, m *protocol.Message) {
	var (
		err error
		ip  = onlyIp(addr)
	)

	r := &protocol.SetIpRequest{}

	if err = r.ReadFrom(m); nil != err {
		log.Printf("error: from host (%s): %s", ip, err)
		s.sendErrorResponse(addr, protocol.BAD_REQUEST, err.Error())
		return
	}

	preexisting := s.registry.Contains(r.Name)

	if err = s.registry.Put(r.Name, r.Ip); nil != err {
		log.Printf("Host (%s): unable to change IP of (%s) to (%s). %s", ip, r.Name, r.Ip, err)
		s.sendErrorResponse(addr, protocol.INVALID_NAME, err.Error())
		return
	}

	if preexisting {
		log.Printf("Host (%s): changed IP of (%s) to (%s)", ip, r.Name, r.Ip)
	} else {
		log.Printf("Host (%s): initialized IP of (%s) to (%s)", ip, r.Name, r.Ip)
	}

	if err = protocol.SendSetIpResponse(&s.conn, addr, "ok", ""); nil != err {
		log.Printf("Error: failed to send set-ip response to Host (%s). %s", ip, err)
	}
}

func handleSetPublicIpRequest(s *Server, addr net.Addr, m *protocol.Message) {
	var (
		err error
		ip  = onlyIp(addr)
	)

	r := &protocol.SetPublicIpRequest{}

	if err = r.ReadFrom(m); nil != err {
		log.Printf("error: from host (%s): %s", addr.String(), err)
		s.sendErrorResponse(addr, protocol.BAD_REQUEST, err.Error())
		return
	}

	preexisting := s.registry.Contains(r.Name)

	if err = s.registry.Put(r.Name, ip); nil != err {
		log.Printf("Host (%s): unable to change IP of (%s) to (%s). %s", ip, r.Name, ip, err)
		s.sendErrorResponse(addr, protocol.INVALID_NAME, err.Error())
		return
	}

	if preexisting {
		log.Printf("Host (%s): changed IP of (%s) to (%s)", ip, r.Name, ip)
	} else {
		log.Printf("Host (%s): initialized IP of (%s) to (%s)", ip, r.Name, ip)
	}

	if err = protocol.SendSetIpResponse(&s.conn, addr, "ok", ""); nil != err {
		log.Printf("Error: failed to send set-ip response to Host (%s). %s", ip, err)
	}
}
