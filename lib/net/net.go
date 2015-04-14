package net

import (
	"net"
	"sync"
)

type Conn struct {
	net.UDPConn
	wg sync.WaitGroup
}

func (c *Conn) Lend() {
	c.wg.Add(1)
}

func (c *Conn) Close() error {
	c.wg.Wait()
	return c.UDPConn.Close()
}

func (c *Conn) Release() {
	c.wg.Done()
}
