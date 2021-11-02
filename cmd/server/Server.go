package main

import (
	"Samurai/net"
	"log"
)

type PingRouter struct {
	net.BaseRouter
}

func (p *PingRouter) PreHandle(req net.Req) {
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("Before!"))
	if err != nil {
		log.Fatalln(err)
	}
}

func (p *PingRouter) Handle(req net.Req) {
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("Ping!"))
	if err != nil {
		log.Fatalln(err)
	}
}

func (p *PingRouter) PostHandle(req net.Req) {
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("After!"))
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	s := net.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
