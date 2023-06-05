package iface

type IHandler interface {
	PrevHandle(c IContext)
	LastHandle(c IContext)
}
