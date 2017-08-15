package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Server struct {
	Temp uint32
	Mux  *http.ServeMux
}

func (s *Server) HandleFoo() {
	s.Mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world\n")
		s.Temp = 12
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	s := &Server{
		Temp: 1,
		Mux:  http.NewServeMux(),
	}
	s.HandleFoo()
	go http.ListenAndServe(":8875", s.Mux)
	time.Sleep(10 * time.Second)
	log.Println(s.Temp)

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
