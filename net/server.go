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
	AddRouter(r router)
}

type Server struct {
	Name    string
	Version string
	IP      string
	Port    int
	Router  router
}

func (s *Server) Start() {
	log.Printf("%s start listen at %s:%d\n", s.Name, s.IP, s.Port)
	log.Printf("Version : %s, MaxConn : %d, MaxPacketSize : %d",
		GlobalConfig.Version, GlobalConfig.MaxConn, GlobalConfig.MaxPacketSize)
	go func() {
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
				log.Fatalln("Accept err: ", err)
				continue
			}
			// todo : set max connection
			// todo : handle new connection function

			dealConnection := NewConnection(conn, connID, s.Router)
			connID++

			go dealConnection.Start()
		}
	}()
}

func (s *Server) Stop() {
	log.Println("Stop server ", s.Name)
	// todo : clean the connection
}

func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(r router) {
	s.Router = r
}

func NewServer() server {
	GlobalConfig.Reload()

	return &Server{
		Name:    GlobalConfig.Name,
		Version: "tcp4",
		IP:      GlobalConfig.Host,
		Port:    GlobalConfig.Port,
		Router:  nil,
	}
}
