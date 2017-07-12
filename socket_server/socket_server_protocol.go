package socket_server

import (
	"errors"
	"github.com/gansidui/gotcp"
)

var (
	ErrPeerClosed = errors.New("peer closed")
)

type Raw struct {
	raw []byte
}

func (r *Raw) Serialize() {
	return r.raw
}

func (ss *SocketServer) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	data := make([]byte, 1024)
	length, err := conn.Read(data)
	log.Printf("<IN> %x  %x\n", conn, data[0:length])
	if err != nil {
		return nil, err
	}

	if readLengh == 0 {
		return nil, ErrPeerClosed
	}

	return &Raw{
		raw: data[0:length],
	}, nil
}
