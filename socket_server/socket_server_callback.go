package socket_server

import (
	"github.com/gansidui/gotcp"
	"github.com/giskook/bed2/socket_server/protocol"
)

func (ss *SocketServer) OnConnect(c *gotcp.Conn) bool {
	connection := NewConnection(c, &ConnConf{
		read_limit:  ss.conf.ReadLimit,
		write_limit: ss.conf.WriteLimit,
		block_size:  ss.conf.BlockSize,
	})

	c.PutExtraData(connection)
	go connection.Check()

	return true
}

func (ss *SocketServer) OnClose(c *gotcp.Conn) {
	connection := c.GetExtraData().(*Connection)
	if connection == ss.cm.Get(connection.ID) {
		ss.cm.Del(connection.ID)
	}
	connection.Close()
}

func (ss *SocketServer) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	connection := c.GetExtraData().(*Connection)
	connection.UpdateReadFlag()
	connection.RecvBuffer.Write(p.Serialize())
	for {
		protocol_id, length := protocol.CheckProtocol(connection.RecvBuffer)
		buf := make([]byte, length)
		connection.RecvBuffer.Read(buf)
		switch protocol_id {
		case protocol.PROTOCOL_HALF_PACK:
			return true
		case protocol.PROTOCOL_ILLEGAL:
			return true
		case protocol.PROTOCOL_REQ_LOGIN:
			ss.eh_login(buf, connection)
		case protocol.PROTOCOL_REQ_HEART:
			ss.eh_heart(buf, connection)
		case protocol.PROTOCOL_REP_CONTROL:
			ss.eh_control(buf)
		case protocol.PROTOCOL_REP_ACTIVE_TEST:
			ss.eh_active_test(buf)
		case protocol.PROTOCOL_REP_TRANSPARENT_TRANSMISSION:
			ss.eh_rep_transparent_transmission(connection, buf)
		case protocol.PROTOCOL_REP_TRANSPARENT_TRANSMISSION_STOP:
			ss.eh_rep_transparent_transmission_stop(connection, buf)
		case protocol.PROTOCOL_REQ_BIN_FILE:
			ss.eh_req_bin_file(connection, buf)
		}
	}
}
