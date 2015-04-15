package buffer


const (
	MAX_UDP_PACKET_SIZE = 65535
)

func CreateUdpBuffer() []byte {
	return make([]byte, MAX_UDP_PACKET_SIZE)
}

