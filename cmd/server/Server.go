package main

import (
	"Samurai/net"
	"log"
)

type PingRouter struct {
	net.BaseRouter
}

func (p *PingRouter) Handle(req net.Req) {
	log.Printf("Receive from client : msgID=%d,data=%s\n", req.GetMsgID(), string(req.GetData()))
	err := req.GetConnection().SendMsg(1, []byte("Ping!"))
	if err != nil {
		log.Fatalln(err)
	}
}

type HelloRouter struct {
	net.BaseRouter
}

func (h *HelloRouter) Handle(req net.Req) {
	log.Printf("Receive from client : msgID=%d,data=%s\n", req.GetMsgID(), string(req.GetData()))
	err := req.GetConnection().SendMsg(1, []byte("Hi!"))
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	s := net.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
