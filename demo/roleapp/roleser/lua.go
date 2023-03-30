package roleser

import (
	lua "github.com/yuin/gopher-lua"
)

// 调起Lua，校验玩家是否允许登陆
func LuaRoleHook(luafile string, roleid string) (bool, error) {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(luafile); err != nil {
		return false, err
	}
	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("RoleHook"),
		NRet:    1,
		Protect: true,
	}, lua.LString(roleid))
	if err != nil {
		return false, err
	}
	ret := L.Get(-1)
	L.Pop(1)
	if res, ok := ret.(lua.LBool); ok {
		return bool(res), nil
	} else {
		return false, nil
	}
}
