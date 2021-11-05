package net

import (
	"errors"
	"log"
	"sync"
)

type manager interface {
	Add(conn Conn)
	Remove(conn Conn)
	Get(connID uint32) (Conn, error)
	Len() int
	Clear()
}

type ConnManager struct {
	conns map[uint32]Conn
	lock  sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: make(map[uint32]Conn),
	}
}

func (c *ConnManager) Add(conn Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.conns[conn.GetConnID()] = conn

	log.Printf("Add Conn! All Conn : %d", c.Len())
}

func (c *ConnManager) Remove(conn Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.conns, conn.GetConnID())
	log.Printf("Remove connID=%d", conn.GetConnID())
}

func (c *ConnManager) Get(connID uint32) (Conn, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if conn, ok := c.conns[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Connection not found ")
	}
}

func (c *ConnManager) Len() int {
	return len(c.conns)
}

func (c *ConnManager) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for connID, conn := range c.conns {
		conn.Stop()
		delete(c.conns, connID)
	}
	log.Println("Close all connections.")
}
