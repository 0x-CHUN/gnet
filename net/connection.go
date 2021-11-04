package net

import (
	"errors"
	"io"
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
	// SendMsg send message
	SendMsg(msgID uint32, data []byte) error
}

type HandlerFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	MsgHandler msgHandle
	ExitChan   chan bool // channel for notify
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler msgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		pack := NewPacket()

		header := make([]byte, pack.GetHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), header); err != nil {
			log.Println("Read header error : ", err)
			c.ExitChan <- true
			continue
		}
		msg, err := pack.Unpack(header)
		if err != nil {
			log.Println("Unpack error : ", err)
			c.ExitChan <- true
			continue
		}
		var data []byte
		if msg.GetLen() > 0 {
			data = make([]byte, msg.GetLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Println("Unpack error : ", err)
				c.ExitChan <- true
				continue
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(&req)
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection is closed when send msg. ")
	}
	pack := NewPacket()
	msg, err := pack.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		log.Printf("%d pack error\n", msgID)
		return err
	}
	if _, err := c.Conn.Write(msg); err != nil {
		log.Printf("%d write msg error", msgID)
		c.ExitChan <- true
		return err
	}
	return nil
}
