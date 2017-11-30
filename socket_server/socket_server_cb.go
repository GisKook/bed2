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

func (ss *SocketServer) eh_rep_transparent_transmission(c *Connection, p []byte) {
	ss.SocketOutResult <- protocol.ParseRepTransparentTransmission(p)
	//go c.GoTT()
		c.SetMode(CONNECTION_TRANSPARENT_TRANSMISSION_MODE)
}

func (ss *SocketServer) eh_rep_transparent_transmission_stop(c *Connection, p []byte) {
	//c.StopTT()
	c.SetMode(CONNECTION_LOGIN)
}

func (ss *SocketServer) eh_req_bin_file(c *Connection, p []byte) {
	bin := protocol.ParseReqBin(p)
	c.SendBinBytes(int(bin.Cursor), int(bin.BlockSize))
}
