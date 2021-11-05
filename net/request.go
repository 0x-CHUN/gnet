package net

type Req interface {
	GetConnection() Conn
	GetData() []byte
	GetMsgID() uint32
}

type Request struct {
	conn Conn
	msg  message
}

func (r *Request) GetConnection() Conn {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
