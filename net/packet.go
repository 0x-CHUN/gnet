package net

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type packet interface {
	GetHeaderLen() uint32
	Pack(msg message) ([]byte, error)
	Unpack([]byte) (message, error)
}

type Packet struct {
}

func NewPacket() *Packet {
	return &Packet{}
}

func (p *Packet) GetHeaderLen() uint32 {
	// ID uint32(4B) + Len uint32(4B)
	return 8
}

func (p *Packet) Pack(msg message) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	// write len
	if err := binary.Write(buf, binary.LittleEndian, msg.GetLen()); err != nil {
		return nil, err
	}
	// write message id
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	// write data
	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Packet) Unpack(data []byte) (message, error) {
	reader := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(reader, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	if GlobalConfig.MaxPacketSize > 0 && msg.Len > GlobalConfig.MaxPacketSize {
		return nil, errors.New("Too large message. ")
	}
	return msg, nil
}
