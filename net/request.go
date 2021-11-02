package net

type Req interface {
	GetConnection() connection
	GetData() []byte
}

type Request struct {
	conn connection
	data []byte
}

func (r *Request) GetConnection() connection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
