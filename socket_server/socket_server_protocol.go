package socket_server

import (
	"errors"
	"github.com/gansidui/gotcp"
	"log"
	"net"
)

var (
	ErrPeerClosed = errors.New("peer closed")
)

type Raw struct {
	raw []byte
}

func (r *Raw) Serialize() []byte {
	return r.raw
}

func (ss *SocketServer) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	data := make([]byte, 1024)
	length, err := conn.Read(data)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, ErrPeerClosed
	}
	log.Printf("<IN>  %x  %x\n", conn, data[0:length])

	return &Raw{
		raw: data[0:length],
	}, nil
}
