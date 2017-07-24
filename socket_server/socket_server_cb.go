package socket_server

import (
	"github.com/giskook/bed2/socket_server/protocol"
)

func (ss *SocketServer) eh_login(p []byte, c *Connection) {
	login := protocol.ParseReqLogin(p)
	c.ID = login.ID
	c.status = CONNECTION_LOGIN
	ss.cm.Put(login.ID, c)
	c.Send(login)
}

func (ss *SocketServer) eh_heart(p []byte, c *Connection) {
	heart := protocol.ParseReqHeart(p)
	if c.status == CONNECTION_LOGIN {
		c.Send(heart)
	}
}

func (ss *SocketServer) eh_control(p []byte) {
	ss.SocketOutResult <- protocol.ParseRepControl(p)
}

func (ss *SocketServer) eh_active_test(p []byte) {
	ss.SocketOutResult <- protocol.ParseRepActive(p)
}
