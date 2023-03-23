package meta

import (
	"fmt"
	"m3game/plugins/log"
	"regexp"

	"github.com/pkg/errors"
)

type RouteType string

func (r RouteType) String() string {
	return string(r)
}

const (
	RouteTypeP2P    RouteType = "RouteTypeP2P"
	RouteTypeRandom RouteType = "RouteTypeRandom"
	RouteTypeHash   RouteType = "RouteTypeHash"
	RouteTypeBroad  RouteType = "RouteTypeBroad"
	RouteTypeMulti  RouteType = "RouteTypeMulti"
	RouteTypeSingle RouteType = "RouteTypeSingle"
)

type IsNty string

func (r IsNty) String() string {
	return string(r)
}

const (
	IsNtyTrue  IsNty = "true"
	IsNtyFalse IsNty = "false"
)

var (
	regexAppID   *regexp.Regexp
	regexSvcID   *regexp.Regexp
	regexWorldID *regexp.Regexp
	regexEnvID   *regexp.Regexp
)

func init() {
	var err error
	if regexAppID, err = regexp.Compile("^([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)$"); err != nil {
		log.Fatal("regexAppID.Compile err %s", err)
	}
	if regexSvcID, err = regexp.Compile("^([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)$"); err != nil {
		log.Fatal("regexSvcID.Compile err %s", err)
	}
	if regexWorldID, err = regexp.Compile("^([^\\.]+)\\.([^\\.]+)$"); err != nil {
		log.Fatal("regexWorldID.Compile err %s", err)
	}
	if regexEnvID, err = regexp.Compile("^([^\\.]+)$"); err != nil {
		log.Fatal("regexEnvID.Compile err %s", err)
	}

}

type RouteApp string

func (r RouteApp) String() string {
	return string(r)
}

func (r RouteApp) Vaild() bool {
	groups := regexAppID.FindStringSubmatch(string(r))
	if len(groups) == 0 {
		return false
	}
	return true
}

func (r RouteApp) Parse() (env string, world string, fun string, ins string, err error) {
	err = nil
	groups := regexAppID.FindStringSubmatch(string(r))
	if len(groups) == 0 {
		err = fmt.Errorf("RouteApp Parse fail %s", string(r))
	} else {
		env = groups[1]
		world = groups[2]
		fun = groups[3]
		ins = groups[4]
	}
	return
}

func GenRouteApp(env string, world string, fun string, ins string) RouteApp {
	return RouteApp(fmt.Sprintf("%s.%s.%s.%s", env, world, fun, ins))
}

type RouteSvc string

func (r RouteSvc) String() string {
	return string(r)
}
func (r RouteSvc) Vaild() bool {
	groups := regexSvcID.FindStringSubmatch(string(r))
	if len(groups) == 0 {
		return false
	}
	return true
}

func (r RouteSvc) Parse() (env string, world string, fun string, err error) {
	err = nil
	groups := regexSvcID.FindStringSubmatch(string(r))
	if len(groups) == 0 {
		err = fmt.Errorf("RouteSvc Parse fail %s", string(r))
	} else {
		env = groups[1]
		world = groups[2]
		fun = groups[3]
	}
	return
}

func GenRouteSvc(env string, world string, fun string) RouteSvc {
	return RouteSvc(fmt.Sprintf("%s.%s.%s", env, world, fun))
}

type RouteWorld string

func (r RouteWorld) String() string {
	return string(r)
}
func (r RouteWorld) Vaild() bool {
	groups := regexWorldID.FindStringSubmatch(string(r))
	if len(groups) == 0 {
		return false
	}
	return true
}

func (r RouteWorld) Parse() (env string, world string, err error) {
	err = nil
	groups := regexWorldID.FindStringSubmatch(string(r))
	if len(groups) == 0 {
		err = fmt.Errorf("RouteWorld Parse fail %s", string(r))
	} else {
		env = groups[1]
		world = groups[2]
	}
	return
}

func GenRouteWorld(env string, world string) RouteWorld {
	return RouteWorld(fmt.Sprintf("%s.%s", env, world))
}

type RouteEnv string

func (r RouteEnv) String() string {
	return string(r)
}
func (r RouteEnv) Vaild() bool {
	if len(r) == 0 {
		return false
	}
	return true
}

func (r RouteEnv) Parse() (env string, world string, err error) {
	if len(r) == 0 {
		err = errors.New("RouteEnv Parse fail")
		return
	}
	env = string(r)
	return
}

func GenRouteEnv(env string) RouteEnv {
	return RouteEnv(env)
}
