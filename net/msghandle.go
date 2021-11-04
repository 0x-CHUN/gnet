package net

import "log"

type msgHandle interface {
	DoMsgHandler(req Req)
	AddRouter(msgID uint32, r router)
}

type MsgHandle struct {
	Apis map[uint32]router
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]router),
	}
}

func (m *MsgHandle) DoMsgHandler(req Req) {
	handler, ok := m.Apis[req.GetMsgID()]
	if !ok {
		log.Printf("Api msgID=%d is not found!", req.GetMsgID())
		return
	}
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (m *MsgHandle) AddRouter(msgID uint32, r router) {
	if _, ok := m.Apis[msgID]; ok {
		log.Printf("Repeated api , msgID=%d", msgID)
		return
	}
	m.Apis[msgID] = r
	log.Printf("Add api msgID=%d", msgID)
}
