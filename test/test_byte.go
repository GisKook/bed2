package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var HTTP_REQUEST_DESC []string = []string{
	"成功",
	"缺少参数",
	"服务器内部错误"}

type clint int

func (code clint) Color() {
	log.Println(HTTP_REQUEST_DESC[code])
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//var test clint
	//test = 1
	//test.Color()
	for _, v := range HTTP_REQUEST_DESC {
		log.Println(v)
	}
	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
