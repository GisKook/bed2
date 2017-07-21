package socket_server

import (
	"github.com/gansidui/gotcp"
	"github.com/giskook/bed2/base"
	"github.com/giskook/bed2/conf"
	"log"
	"net"
	"time"
)

type SocketServer struct {
	conf            *conf.SocketConf
	srv             *gotcp.Server
	cm              *ConnMgr
	SocketOutResult chan *base.SocketResponse
}

func NewSocketServer(conf *conf.SocketConf) *SocketServer {
	return &SocketServer{
		conf:            conf,
		cm:              NewConnMgr(),
		SocketOutResult: make(chan *base.SocketResponse, 1024),
	}
}

func (ss *SocketServer) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+ss.conf.BindPort)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}

	ss.srv = gotcp.NewServer(config, ss, ss)

	go ss.srv.Start(listener, time.Second)
	log.Println("socket Listening:", listener.Addr())

	return nil
}

func (ss *SocketServer) Send(id uint64, p gotcp.Packet) error {
	c := ss.cm.Get(id)
	if c != nil {
		return c.Send(p)
	}

	return base.ErrBedOffline

}

func (ss *SocketServer) Close() {
	close(ss.SocketOutResult)
	ss.srv.Stop()
}
