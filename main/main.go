package main

import (
	"fmt"
	"github.com/giskook/bed2/conf"
	"github.com/giskook/bed2/reactor"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := conf.ReadConfig("./conf.json")
	checkError(err)

	reactor := reactor.NewReactor(configuration)
	reactor.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	reactor.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
