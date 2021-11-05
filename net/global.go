package net

import (
	"encoding/json"
	"io/ioutil"
)

type Global struct {
	TCPServer      server
	Host           string
	Port           int
	Name           string
	Version        string
	MaxPacketSize  uint32
	MaxConn        int
	WorkerPoolSize uint32
	MaxWorkTask    uint32
	MaxBufferSize  uint32
	ConfFilePath   string
}

var GlobalConfig *Global

func (g *Global) Reload() {
	data, err := ioutil.ReadFile(g.ConfFilePath)
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
		Name:           "Global",
		Version:        "version 1.0",
		Port:           9999,
		Host:           "0.0.0.0",
		MaxConn:        12000,
		MaxPacketSize:  4096,
		WorkerPoolSize: 10,
		MaxWorkTask:    1024,
		MaxBufferSize:  1024,
		ConfFilePath:   "config/config.json",
	}
	GlobalConfig.Reload()
}
