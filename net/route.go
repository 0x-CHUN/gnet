package net

type router interface {
	PreHandle(req Req)
	Handle(req Req)
	PostHandle(req Req)
}

type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(_ Req) {
}

func (b *BaseRouter) Handle(_ Req) {
}

func (b *BaseRouter) PostHandle(_ Req) {

}
