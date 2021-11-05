package net

import (
	"fmt"
	"log"
	"net"
	"time"
)

type server interface {
	// Start server
	Start()
	// Stop the server
	Stop()
	// Serve start to serve
	Serve()
	// AddRouter add router to server
	AddRouter(msgID uint32, r router)
	// GetManager get the Conn manager
	GetManager() manager
	SetOnStart(func(conn Conn))
	SetOnStop(func(conn Conn))
	CallOnStart(conn Conn)
	CallOnStop(conn Conn)
}

type Server struct {
	Name       string
	Version    string
	IP         string
	Port       int
	MsgHandler msgHandle
	Manager    manager
	OnStart    func(conn Conn)
	OnStop     func(conn Conn)
}

func (s *Server) Start() {
	log.Printf("%s start listen at %s:%d\n", s.Name, s.IP, s.Port)
	log.Printf("Version : %s, MaxConn : %d, MaxPacketSize : %d",
		GlobalConfig.Version, GlobalConfig.MaxConn, GlobalConfig.MaxPacketSize)
	go func() {
		// start worker pool
		s.MsgHandler.StartWorkerPool()

		// get a tcp addr
		addr, err := net.ResolveTCPAddr(s.Version, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Fatalln("Resolve tcp addr err:", err)
			return
		}

		// listen to the addr
		listener, err := net.ListenTCP(s.Version, addr)
		if err != nil {
			log.Fatalf("Listen %s, err : %t", s.Version, err)
			return
		}

		// todo : id generator
		var connID uint32
		connID = 0

		// start to listen
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Println("Accept err: ", err)
				continue
			}
			// too much Conn,drop the new Conn
			if s.Manager.Len() >= GlobalConfig.MaxConn {
				log.Println("Too much Conn, close the Conn from ", conn.RemoteAddr())
				err := conn.Close()
				if err != nil {
					log.Println("Close connection err ", err)
				}
				continue
			}

			dealConnection := NewConnection(s, conn, connID, s.MsgHandler)
			connID++

			go dealConnection.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.Manager.Clear()
	log.Println("Stop server ", s.Name)
}

func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgID uint32, r router) {
	s.MsgHandler.AddRouter(msgID, r)
}

func (s *Server) GetManager() manager {
	return s.Manager
}

func (s *Server) SetOnStart(f func(conn Conn)) {
	s.OnStart = f
}

func (s *Server) SetOnStop(f func(conn Conn)) {
	s.OnStop = f
}

func (s *Server) CallOnStart(conn Conn) {
	if s.OnStart != nil {
		s.OnStart(conn)
	}
}

func (s *Server) CallOnStop(conn Conn) {
	if s.OnStop != nil {
		s.OnStop(conn)
	}
}

func NewServer() server {
	GlobalConfig.Reload()

	return &Server{
		Name:       GlobalConfig.Name,
		Version:    "tcp4",
		IP:         GlobalConfig.Host,
		Port:       GlobalConfig.Port,
		MsgHandler: NewMsgHandle(),
		Manager:    NewConnManager(),
	}
}
