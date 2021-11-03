package net

type Req interface {
	GetConnection() connection
	GetData() []byte
	GetMsgID() uint32
}

type Request struct {
	conn connection
	msg  message
}

func (r *Request) GetConnection() connection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
