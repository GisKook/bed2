package socket_server

import (
	"github.com/giskook/bed2/socket_server/protocol"
)

func (ss *SocketServer) OnConnect(c *gotcp.Conn) bool {
	connection := NewConnection(c, &ConnConf{
		read_limit:  ss.conf.ReadLimit,
		write_limit: ss.conf.WriteLimit,
	})

	c.PutExtraData(connection)
	go c.Check()
}

func (ss *SocketServer) OnClose(c *gotcp.Conn) {
	connection := c.GetExtraData().(*Connection)
	ss.cm.Del(connection.ID)
	connection.Close()
}

func (ss *SocketServer) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	connection := c.GetExtraData().(*Connection)
	connection.RecvBuffer.Write(p.Serialize())
	for {
		protocol_id, length := protocol.CheckProtocol(connection.RecvBuffer)
		switch protocol_id {
		case protocol.PROTOCOL_HALF_PACK:
			return true
		case protocol.PROTOCOL_ILLEGAL:
			return true
		case protocol.PROTOCOL_REQ_LOGIN:
		case protocol.PROTOCOL_REQ_HEART:
		case protocol.PROTOCOL_REP_CONTROL:
		case protocol.PROTOCOL_REP_ACTIVE_TEST:
		}
	}
}
