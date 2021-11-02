package net

import (
	"encoding/json"
	"io/ioutil"
)

type Global struct {
	TCPServer     server
	Host          string
	Port          int
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       int
}

var GlobalConfig *Global

func (g *Global) Reload() {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalConfig = &Global{
		Name:          "Global",
		Version:       "version 1.0",
		Port:          9999,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}
	GlobalConfig.Reload()
}
