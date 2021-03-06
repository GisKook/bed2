package http_server

import (
	"github.com/giskook/bed2/base"
	"github.com/giskook/bed2/conf"
	"log"
	"net/http"
)

const (
	HTTP_ACTIVE_TEST              string = "/olderhc/bed/checkOnline"
	HTTP_CONTROL                  string = "/olderhc/bed/control"
	HTTP_TRANSPARENT_TRANSMISSION string = "/olderhc/bed/transparent_transmission"
	HTTP_LOAD_BIN                 string = "/olderhc/bed/load_bin"
)

type HttpServer struct {
	HttpRequestChan chan *base.HttpRequest
	router          *HttpRouter
	conf            *conf.HttpConf
	mux             *http.ServeMux
}

func NewHttpServer(conf *conf.HttpConf) *HttpServer {
	return &HttpServer{
		HttpRequestChan: make(chan *base.HttpRequest),
		router:          NewHttpRouter(),
		conf:            conf,
		mux:             http.NewServeMux(),
	}
}

func (h *HttpServer) Close() {
	close(h.HttpRequestChan)
}

func (h *HttpServer) Response(resp *base.SocketResponse) {
	h.router.DoResponse(resp)
}

func (h *HttpServer) request(req *base.HttpRequest) {
	h.HttpRequestChan <- req
}

func (h *HttpServer) Init() {
	h.HandleActiveTest(HTTP_ACTIVE_TEST)
	h.HandleControl(HTTP_CONTROL)
	h.HandleTransparentTransmission(HTTP_TRANSPARENT_TRANSMISSION)
	h.HandleLoadBin(HTTP_LOAD_BIN)
}

func (h *HttpServer) Start() error {
	defer func() {
		log.Println("http   Listening", h.conf.Addr)
	}()
	return http.ListenAndServe(h.conf.Addr, h.mux)
}

func (h *HttpServer) makeRequest() (chan *base.SocketResponse, uint32) {
	return h.router.makeRequest()
}
