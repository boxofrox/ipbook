package server

import (
	"net"
	"strings"
)

func onlyIp(addr net.Addr) string {
	ip := addr.String()

	if i := strings.LastIndex(ip, ":"); -1 < i {
		return ip[0:i]
	}

	return ip
}
