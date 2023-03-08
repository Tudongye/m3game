package router

var (
	_router Router
)

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
	GetAllInstances(svcid string) ([]Instance, error)
}

func Set(r Router) {
	if _router != nil {
		panic("Router Only One")
	}
	_router = r
}

func Get() Router {
	if _router == nil {
		panic("Router Mush Have One")
	}
	return _router
}

func Register(app AppReciver) error {
	return Get().Register(app)
}

func Deregister(app AppReciver) error {
	return Get().Deregister(app)
}

func GetAllInstances(svcid string) ([]Instance, error) {
	return Get().GetAllInstances(svcid)
}
