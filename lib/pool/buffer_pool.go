package pool

type BufferPool struct {
	pool   chan []byte
	create func() []byte
}

type BufferAllocator func() []byte

func New(size int, createBuffer BufferAllocator) *BufferPool {
	return &BufferPool{
		pool:   make(chan []byte, size),
		create: createBuffer,
	}
}

func (p *BufferPool) GetFreeBuffer() []byte {
	select {
	case buffer := <-p.pool:
		return buffer

	default:
		return p.create()
	}
}

func (p *BufferPool) Recycle(buffer []byte) {
	if cap(p.pool) < len(p.pool) {
		p.pool <- buffer
	}
}
