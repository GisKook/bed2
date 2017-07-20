package http_server

import (
	"github.com/giskook/bed2/base"
	"sync"
	"sync/atomic"
)

type HttpRouter struct {
	request map[uint32]chan *base.SocketResponse
	index   uint32
	mutex   *sync.RWMutex
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{
		request: make(map[uint32]chan *base.SocketResponse),
		mutex:   new(sync.RWMutex),
	}
}

func (h *HttpRouter) makeRequest() (chan *base.SocketResponse, uint32) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	index := atomic.AddUint32(&h.index, 1)
	rc := make(chan *base.SocketResponse)
	h.request[index] = rc

	return rc, index
}

func (h *HttpRouter) DoResponse(resp *base.SocketResponse) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if c, ok := h.request[resp.SerialID]; ok {
		c <- resp
		delete(h.request, resp.SerialID)
	}
}

func (h *HttpRouter) delRequest(index uint32) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.request, index)
}
