package conf

import (
	"encoding/json"
	"os"
)

type HttpConf struct {
	Addr    string
	TimeOut int
}

type SocketConf struct {
	BindPort      string
	CheckInterval int
	ReadLimit     int
	WriteLimit    int
	ConnTimeOut   int
	SendChanLimit int
	RecvChanLimit int
}
type Configuration struct {
	Http   *HttpConf
	Socket *SocketConf
}

var G_conf *Configuration

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	G_conf = &config

	return &config, err
}

func GetConf() *Configuration {
	return G_conf
}
