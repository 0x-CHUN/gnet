package net

type message interface {
	GetLen() uint32
	GetMsgID() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetData([]byte)
	SetLen(uint322 uint32)
}

type Message struct {
	ID   uint32
	Len  uint32
	Data []byte
}

func NewMsgPacket(id uint32, data []byte) *Message {
	return &Message{
		ID:   id,
		Len:  uint32(len(data)),
		Data: data,
	}
}

func (m *Message) GetLen() uint32 {
	return m.Len
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetLen(len uint32) {
	m.Len = len
}

func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
