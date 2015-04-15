package pool

import "testing"

func createBuffer() []byte { return make([]byte, 1) }

func spyAllocator(flag *bool) BufferAllocator {
	return func() []byte {
		*flag = !*flag
		return createBuffer()
	}
}

func Test_GetFreeBuffer_BeforeAnyRecycle_CreatesBuffer(t *testing.T) {
	allocatorCalled := false

	pool := New(1, spyAllocator(&allocatorCalled))

	pool.GetFreeBuffer()

	if false == allocatorCalled {
		t.Errorf("buffer pool did not run create buffer function")
	}
}

func Test_GetFreeBuffer_AfterAnyRecycle_BufferFromPool(t *testing.T) {
	allocatorNotCalled := true

	pool := New(1, spyAllocator(&allocatorNotCalled))

	// first create and recycle a buffer
	buffer := pool.GetFreeBuffer()
	pool.Recycle(buffer)

	// reset flag
	allocatorNotCalled = true

	// run test
	pool.GetFreeBuffer()

	// assert flag did not change
	if false == allocatorNotCalled {
		t.Errorf("buffer pool did not recycle the buffer")
	}
}

func Test_New_NegativeOrZeroSize_Panics(t *testing.T) {
	t.Skipf("to be written")
}
