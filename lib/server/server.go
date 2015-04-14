package server

import (
	"log"
	"net"
	"sync"

	boxnet "github.com/boxofrox/ipbook/lib/net"
	"github.com/boxofrox/ipbook/lib/pool"
	"github.com/boxofrox/ipbook/lib/protocol"
	"github.com/boxofrox/ipbook/lib/registry"
)

type Server struct {
	conn     boxnet.Conn
	pool     *pool.BufferPool
	registry *registry.Registry
	once     *sync.Once
	done     bool
}

func New(port int) (*Server, error) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}

	c, err := net.ListenUDP("udp", &addr)
	if nil != err {
		return nil, err
	}

	s := Server{
		conn:     boxnet.Conn{UDPConn: *c},
		pool:     pool.New(5, createBuffer),
		registry: registry.New(),
		once:     &sync.Once{},
	}

	return &s, nil
}

func (s *Server) asyncHandleRequest(addr *net.UDPAddr, n int, buffer []byte) {
	var err error

	s.conn.Lend()
	defer s.conn.Release()

	defer s.pool.Recycle(buffer)

	object, err := protocol.Decode(buffer[0:n])
	if nil != err {
		log.Printf("Error: unable to decode request. %s", err)
		protocol.SendErrorResponse(&s.conn, addr, protocol.BAD_REQUEST, "unable to decode request")
		return
	}

	handler, exists := handlers[object.GetType()]
	if !exists {
		log.Printf("Error: unknown request from Host (%s).\n  %s", addr.String(), buffer[0:n])
		return
	}

	handler(s, addr, object)
}

func (s *Server) Listen() {
	// only run the server once per instance
	defer s.reset()
	s.once.Do(func() { s.listen() })
}

func (s *Server) listen() {
	s.done = false

	defer s.conn.Close()

	for !s.done {
		buffer := s.pool.GetFreeBuffer()
		n, addr, err := s.conn.ReadFromUDP(buffer)

		if err != nil {
			log.Printf("Error: reading udp packet from %s. %s", addr.String(), err)

			if nil != addr {
				protocol.SendErrorResponse(&s.conn, addr, protocol.BAD_REQUEST, "unable to read request")
			}

			continue
		}

		go s.asyncHandleRequest(addr, n, buffer)
	}
}

func (s *Server) reset() {
	s.once = &sync.Once{}
}

func (s *Server) Stop() {
	s.done = true
	s.conn.Close()
}

const (
	MAX_UDP_PACKET_SIZE = 65535
)

func createBuffer() []byte {
	return make([]byte, MAX_UDP_PACKET_SIZE)
}

type RequestHandler func(s *Server, addr *net.UDPAddr, msg protocol.Messager)

var handlers = map[int]RequestHandler{
	protocol.TYPE_GET_IP_REQUEST: handleGetIpRequest,
	protocol.TYPE_SET_IP_REQUEST: handleSetIpRequest,
}

func handleGetIpRequest(s *Server, addr *net.UDPAddr, msg protocol.Messager) {
	var err error

	request := msg.(*protocol.GetIpRequest)

	if false == s.registry.Contains(request.Name) {

		log.Printf("Host (%s): requested name (%s) not found in registry.", addr.String(), request.Name)
		if err = protocol.SendErrorResponse(&s.conn, addr, protocol.NAME_NOT_FOUND, "name not found in registry"); nil != err {
			log.Printf("Error: failed to send error response to Host (%s). %s", addr.String(), err)
		}
		return
	}

	ip, _ := s.registry.Get(request.Name)

	if err = protocol.SendGetIpResponse(&s.conn, addr, request.Name, ip); nil != err {
		log.Printf("Error: unable to send get-ip response to Host (%s). %s", addr.String(), err)
		return
	}

	log.Printf("Host (%s): requested IP of (%s).", addr.String(), request.Name)
}

func handleSetIpRequest(s *Server, addr *net.UDPAddr, msg protocol.Messager) {
	var err error

	request := msg.(*protocol.SetIpRequest)
	preexisting := s.registry.Contains(request.Name)

	if err = s.registry.Put(request.Name, request.Ip); nil != err {
		log.Printf("Host (%s): unable to change IP of (%s) to (%s). %s", addr.String(), request.Name, request.Ip, err)

		if err = protocol.SendErrorResponse(&s.conn, addr, protocol.INVALID_NAME, err.Error()); nil != err {
			log.Printf("Error: failed to send error response to Host (%s). %s", addr.String(), err)
		}

		return
	}

	if preexisting {
		log.Printf("Host (%s): changed IP of (%s) to (%s)", addr.String(), request.Name, request.Ip)
	} else {
		log.Printf("Host (%s): initialized IP of (%s) to (%s)", addr.String(), request.Name, request.Ip)
	}

	if err = protocol.SendSetIpResponse(&s.conn, addr, "ok", ""); nil != err {
		log.Printf("Error: failed to send set-ip response to Host (%s). %s", addr.String(), err)
	}
}
