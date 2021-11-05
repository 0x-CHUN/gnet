package net

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type Conn interface {
	// Start the Conn,let the Conn work
	Start()
	// Stop the Conn
	Stop()
	// GetTCPConnection get the original socket tcp Conn
	GetTCPConnection() *net.TCPConn
	// GetConnID get the Conn ID
	GetConnID() uint32
	// RemoteAddr get the remote address
	RemoteAddr() net.Addr
	// SendMsg send message
	SendMsg(msgID uint32, data []byte) error
	// SendBuffer send message to buffer
	SendBuffer(msgID uint32, data []byte) error
	// SetProperty set the property
	SetProperty(key string, val interface{})
	// GetProperty get the property
	GetProperty(key string) (interface{}, error)
	// RemoveProperty Remove the property
	RemoveProperty(key string)
}

type HandlerFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Server      server
	Conn        *net.TCPConn
	ConnID      uint32
	isClosed    bool
	MsgHandler  msgHandle
	ExitChan    chan bool   // channel for notify
	MsgChan     chan []byte // chan for read and write goroutine
	MsgBuffChan chan []byte // buffer chan for read and write goroutine
	Property    map[string]interface{}
	Lock        sync.RWMutex
}

func NewConnection(s server, conn *net.TCPConn, connID uint32, msgHandler msgHandle) *Connection {
	c := &Connection{
		Server:      s,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		ExitChan:    make(chan bool, 1),
		MsgChan:     make(chan []byte),
		MsgBuffChan: make(chan []byte, GlobalConfig.MaxBufferSize),
		Property:    make(map[string]interface{}),
	}
	c.Server.GetManager().Add(c)
	return c
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
		if GlobalConfig.WorkerPoolSize > 0 {
			c.MsgHandler.AddToQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	for {
		select {
		case data := <-c.MsgChan: // some data need to be sent
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("Send data err : ", err)
				return
			}
		case data, ok := <-c.MsgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					log.Println("Send data err : ", err)
					return
				}
			} else {
				log.Println("Message buffer is closed")
				break
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	defer c.Stop()
	// read data goroutine
	go c.StartReader()
	// write data goroutine
	go c.StartWriter()
	// run start hook function
	c.Server.CallOnStart(c)

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
	// run stop hook function
	c.Server.CallOnStop(c)

	err := c.Conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	c.Server.GetManager().Remove(c)

	c.ExitChan <- true
	// close all channel
	close(c.ExitChan)
	close(c.MsgChan)
	close(c.MsgBuffChan)
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
	c.MsgChan <- msg // msg to channel,and send
	return nil
}

func (c *Connection) SendBuffer(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is closed ")
	}
	pack := NewPacket()
	msg, err := pack.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		log.Println("Pack error ", err)
		return err
	}
	c.MsgBuffChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, val interface{}) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Property[key] = val
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	if val, ok := c.Property[key]; ok {
		return val, nil
	} else {
		return nil, errors.New("Not found ")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Property, key)
}
