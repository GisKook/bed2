package reactor

import (
	"github.com/giskook/bed2/conf"
	"github.com/giskook/bed2/http_server"
	"github.com/giskook/bed2/socket_server"
	"log"
)

type Reactor struct {
	conf *conf.Configuration
	exit chan struct{}

	ss   *socket_server.SocketServer
	http *http_server.HttpServer
}

func NewReactor(conf *conf.Configuration) *Reactor {
	return &Reactor{
		conf: conf,
		exit: make(chan struct{}),
	}
}

func (r *Reactor) Start() {
	r.ss = socket_server.NewSocketServer(r.conf.Socket)
	err := r.ss.Start()
	if err != nil {
		log.Println(err.Error())
	}

	r.http = http_server.NewHttpServer(r.conf.Http)
	r.http.Init()
	go r.http.Start()
	r.shunt()
}

func (r *Reactor) Close() {
	r.ss.Close()
	r.http.Close()
	close(r.exit)
}
