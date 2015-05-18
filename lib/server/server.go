package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	boxnet "github.com/boxofrox/ipbook/lib/net"
	"github.com/boxofrox/ipbook/lib/protocol"
	"github.com/boxofrox/ipbook/lib/registry"
)

type Server struct {
	conn     boxnet.Conn
	registry *registry.Registry
	once     *sync.Once
	done     bool
}

func New(host string, port int) (*Server, error) {
	var (
		addr *net.UDPAddr
		err  error
		c    net.PacketConn
	)

	if addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port)); nil != err {
		return nil, err
	}

	c, err = net.ListenUDP("udp", addr)
	if nil != err {
		return nil, err
	}

	s := Server{
		conn:     boxnet.Conn{PacketConn: c},
		registry: registry.New(),
		once:     &sync.Once{},
	}

	return &s, nil
}

func (s *Server) Run() {
	// only run the server once per instance
	defer s.reset()
	s.once.Do(func() { go s.listen() })

	// terminate gracefully.  ie let server finish responding to requests.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	s.Stop()
}

func (s *Server) Stop() {
	s.done = true
	s.conn.Close()
}

func (s *Server) asyncHandleRequest(pkt *protocol.Packet) {
	var err error

	s.conn.Lend()
	defer s.conn.Release()

	message, err := pkt.ReadMessage()
	if nil != err {
		log.Printf("Error: unable to decode request. %s", err)
		s.sendErrorResponse(pkt.Raddr, protocol.BAD_REQUEST, "unable to decode request")
		return
	}

	handler, exists := handlers[message.Type]
	if !exists {
		log.Printf("Error: unknown request from Host (%s).\n  %s", pkt.Raddr.String(), pkt.Data())
		return
	}

	handler(s, pkt.Raddr, message)
}

func (s *Server) listen() {
	s.done = false

	defer s.conn.Close()

	for !s.done {
		pkt, err := protocol.ReadPacket(&s.conn)

		if err != nil {
			log.Printf("Error: reading udp packet. %s", err)

			if nil != pkt {
				s.sendErrorResponse(pkt.Raddr, protocol.BAD_REQUEST, "unable to read request")
			}

			continue
		}

		go s.asyncHandleRequest(pkt)
	}
}

func (s *Server) reset() {
	s.once = &sync.Once{}
}

func (s *Server) sendErrorResponse(addr net.Addr, code int, reason string) bool {
	if err := protocol.SendErrorResponse(&s.conn, addr, code, reason); nil != err {
		log.Printf("Error: unabled to send error response to %s. %s", addr.String(), err.Error())
		return false
	}
	return true
}

func (s *Server) sendGetIpResponse(addr net.Addr, name, ip string) bool {
	if err := protocol.SendGetIpResponse(&s.conn, addr, name, ip); nil != err {
		log.Printf("Error: unable to send get-ip response to Host (%s). %s", addr.String(), err)
		return false
	}
	return true
}

var handlers = map[int]RequestHandler{
	protocol.TYPE_GET_IP_REQUEST:        handleGetIpRequest,
	protocol.TYPE_SET_IP_REQUEST:        handleSetIpRequest,
	protocol.TYPE_SET_PUBLIC_IP_REQUEST: handleSetPublicIpRequest,
}
