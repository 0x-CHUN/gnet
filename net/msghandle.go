package net

import "log"

type msgHandle interface {
	DoMsgHandler(req Req)
	AddRouter(msgID uint32, r router)
	StartWorkerPool()
	AddToQueue(req Req)
}

type MsgHandle struct {
	Apis           map[uint32]router
	WorkerPoolSize uint32
	TaskQueue      []chan Req
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]router),
		WorkerPoolSize: GlobalConfig.WorkerPoolSize,
		TaskQueue:      make([]chan Req, GlobalConfig.MaxWorkTask),
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

func (m *MsgHandle) StartWorker(workerID int, taskQueue chan Req) {
	log.Printf("Worker ID=%d is working", workerID)
	for { // wait for request
		select {
		case req := <-taskQueue:
			m.DoMsgHandler(req)
		}
	}
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan Req, GlobalConfig.MaxWorkTask)
		go m.StartWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandle) AddToQueue(req Req) {
	workerID := req.GetConnection().GetConnID() % m.WorkerPoolSize
	log.Printf("Add ConnID=%d request msgID=%d to workerID=%d",
		req.GetConnection().GetConnID(), req.GetMsgID(), workerID)
	m.TaskQueue[workerID] <- req
}
