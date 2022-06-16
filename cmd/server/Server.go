package main

import (
	"gnet/net"
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

func OnStart(conn net.Conn) {
	conn.SetProperty("Admin", "Admin")
	log.Println("On Start")
}

func OnClose(conn net.Conn) {
	if val, err := conn.GetProperty("Admin"); err == nil {
		log.Fatalln(val)
	}
	log.Println("On Close")
}

func main() {
	s := net.NewServer()
	s.SetOnStart(OnStart)
	s.SetOnStop(OnClose)
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
