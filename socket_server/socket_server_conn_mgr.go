package socket_server

import (
	"sync"
)

type ConnMgr struct {
	connections map[uint64]*Connection
	mutex       *sync.RWMutex
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connections: make(map[uint64]*Connection),
		mutex:       new(sync.RWMutex),
	}
}

func (cm *ConnMgr) Put(id uint64, c *Connection) {
	cm.mutex.Lock()
	cm.connections[id] = c
	cm.mutex.Unlock()
}

func (cm *ConnMgr) Get(id uint64) *Connection {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if c, ok := cm.connections[id]; ok {
		return c
	}

	return nil
}

func (cm *ConnMgr) Del(id uint64) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.connections, id)
}
