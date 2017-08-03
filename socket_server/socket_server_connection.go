package socket_server

import (
	"bytes"
	"github.com/gansidui/gotcp"
	"github.com/giskook/bed2/base"
	"sync/atomic"
	"time"
)

const (
	CONNECTION_NOT_LOGIN                     uint8 = 0
	CONNECTION_LOGIN                         uint8 = 1
	CONNECTION_TRANSPARENT_TRANSMISSION_MODE uint8 = 2
)

type ConnConf struct {
	read_limit  int
	write_limit int
	block_size  int
}

type Connection struct {
	conf            *ConnConf
	c               *gotcp.Conn
	ID              uint64
	RecvBuffer      *bytes.Buffer
	read_timestamp  int64
	write_timestamp int64
	exit            chan struct{}
	status          uint8

	cursor                             int      // for transparent transmission
	chan_stop_transparent_transmission chan int // for transparent transmission

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
		chan_stop_transparent_transmission: make(chan int),
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

func (c *Connection) Send(p gotcp.Packet) error {
	c.UpdateWriteFlag()
	if c.status != CONNECTION_TRANSPARENT_TRANSMISSION_MODE {
		return c.c.AsyncWritePacket(p, 0)
	}
	return nil
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

func (c *Connection) SetMode(mode uint8) {
	c.status = mode
}

func (c *Connection) SendBinBytes(cursor int, block_size int) {
	c.c.AsyncWritePacket(&Raw{
		raw: base.GetTTB().GetBytes(cursor, block_size),
	}, 0)
}

func (c *Connection) GoTT() {
	c.status = CONNECTION_TRANSPARENT_TRANSMISSION_MODE
	c.cursor = 0
	defer func() {
		c.status = CONNECTION_LOGIN
		base.GetTTB().Decrease()
	}()
	base.GetTTB().Increase()
	block_size := c.conf.block_size
	bin_size := base.GetTTB().GetBinSize()
	for {
		select {
		case <-c.exit:
			return
		case <-c.chan_stop_transparent_transmission:
			return
		default:
			var e error
			///ErrWriteBlocking
			for {
				//time.Sleep(300 * time.Millisecond)
				block_size = calc_block_size(bin_size, c.cursor, block_size)
				e = c.c.AsyncWritePacket(&Raw{
					raw: base.GetTTB().GetBytes(c.cursor, block_size),
				}, 0)
				if e != nil {
					if e != gotcp.ErrWriteBlocking {
						return
					}
				} else {
					c.cursor += block_size
				}
				if c.cursor == bin_size {
					return
				}
			}
		}
	}
}

func calc_block_size(bin_size int, cursor int, block_size int) int {
	if bin_size-cursor < block_size {
		return bin_size - cursor
	}

	return block_size
}

func (c *Connection) StopTT() {
	c.chan_stop_transparent_transmission <- 0
}
