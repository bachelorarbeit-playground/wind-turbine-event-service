package lua

import (
	"fmt"

	luatime "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func RunScript(script []byte, functionName string, input any) (*lua.LTable, error) {
	L := lua.NewState()
	defer L.Close()

	// use => local time = require("time")
	luatime.Preload(L)

	if err := L.DoString(string(script)); err != nil {
		return nil, fmt.Errorf("loading Lua script failed: %s", err.Error())
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal(functionName),
		NRet:    1,
		Protect: true,
	}, luar.New(L, input)); err != nil {
		return nil, fmt.Errorf("executing Lua function '%s' failed: %s", functionName, err.Error())
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.Type() != lua.LTTable {
		return nil, nil
	}

	return result.(*lua.LTable), nil
}
