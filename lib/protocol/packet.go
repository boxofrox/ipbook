package protocol

import (
	"net"
	"runtime"
)

type Packet struct {
	buffer []byte
	size   int
	Raddr  net.Addr
}

func ReadPacket(conn net.PacketConn) (p *Packet, err error) {
	p = &Packet{
		buffer: bufferPool.GetFreeBuffer(),
		size:   0,
		Raddr:  nil,
	}

	runtime.SetFinalizer(p, packetFinalizer)

	if p.size, p.Raddr, err = conn.ReadFrom(p.buffer); nil != err {
		return nil, err
	}

	return p, nil
}

func packetFinalizer(p *Packet) {
	bufferPool.Recycle(p.buffer)
}

func (p *Packet) Data() string {
	return string(p.buffer[0:p.size])
}

func (p *Packet) ReadMessage() (m *Message, err error) {
	m = &Message{}
	err = m.ReadMessage(p.buffer[0:p.size])
	return
}
