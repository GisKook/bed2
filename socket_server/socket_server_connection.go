package socket_server

import (
	"bytes"
	"github.com/gansidui/gotcp"
	"sync/atomic"
	"time"
)

type ConnConf struct {
	read_limit  int
	write_limit int
}

type Connection struct {
	conf            *ConnConf
	c               *gotcp.Conn
	ID              uint64
	RecvBuffer      *bytes.Buffer
	read_timestamp  int64
	write_timestamp int64
	exit            chan struct{}

	ticker *time.Ticker
}

func NewConnection(c *gotcp.Conn, conf *ConnConf) *Connection {
	return &Connection{
		conf:            conf,
		c:               c,
		read_timestamp:  time.Now().Unix(),
		write_timestamp: time.Now().Unix(),
		RecvBuffer:      bytes.NewBuffer([]byte{}),
		ticker:          time.NewTicker(10 * time.Second),
		exit:            make(chan struct{}),
	}
}

func (c *Connection) Close() {
	close(c.exit)
	c.RecvBuffer.Reset()
	c.ticker.Stop()
}

func (c *Connection) UpdateReadFlag() {
	atomic.StoreInt64(&c.read_timestamp, time.Now().Unix())
}

func (c *Connection) UpdateWriteFlag() {
	atomic.StoreInt64(&c.write_timestamp, time.Now().Unix())
}

func (c *Connection) Check() {
	defer func() {
		c.c.Close()
	}()
	for {
		select {
		case <-c.exit:
			return
		case <-c.ticker.C:
			now := time.Now().Unix()
			if now-c.read_timestamp > int64(c.conf.read_limit) ||
				now-c.write_timestamp > int64(c.conf.write_limit) {

				return
			}
		}
	}
}
