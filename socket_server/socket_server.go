package socket_server

import (
	"github.com/gansidui/gotcp"
	"github.com/giskook/conf"
	"net"
	"time"
)

type SocketServer struct {
	conf *conf.SocketConf
	srv  *gotcp.Server
	cm   *ConnMgr
}

func NewSocketServer(conf *conf.SocketConf) *SocketServer {
	return &SocketServer{
		conf: conf,
		cm:   NewConnMgr(),
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
	log.Println("listening:", listener.Addr())
}
