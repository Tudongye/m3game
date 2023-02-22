package router

type AppReciver interface {
	IDStr() string
}

type Instance interface {
	GetHost() string
	GetPort() uint32
	GetIDStr() string
}

type Router interface {
	Register(app AppReciver) error
	Deregister(app AppReciver) error
	GetAllInstances(envid string, worldid string, funcid string) ([]Instance, error)
}
