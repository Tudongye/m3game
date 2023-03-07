package consul

import "m3game/runtime/plugin"

func init() {
	plugin.RegisterFactory(&Factory{})
}
