package net

import (
	"log"
	"net"
)

type connection interface {
	// Start the connection,let the connection work
	Start()
	// Stop the connection
	Stop()
	// GetTCPConnection get the original socket tcp connection
	GetTCPConnection() *net.TCPConn
	// GetConnID get the connection ID
	GetConnID() uint32
	// RemoteAddr get the remote address
	RemoteAddr() net.Addr
}

type HandlerFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	Router   router
	ExitChan chan bool // channel for notify
}

func NewConnection(conn *net.TCPConn, connID uint32, r router) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   r,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			log.Fatalln("Receive err : ", err)
			c.ExitChan <- true
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}
		go func(req Req) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)
	}
}

func (c *Connection) Start() {
	// read data
	go c.StartReader()

	for {
		select {
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// todo : stop function

	err := c.Conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	c.ExitChan <- true
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
