package index

import (
	"fmt"
	"github.com/senyu-up/toolbox/example/internal/script"
)

var scriptRouteMap = map[string]ScriptFunc{
	"init_data":  script.InitData,
	"push_event": script.PushEvent,
}

type ScriptFunc func(map[string]string) error

func GetScriptFunc(scriptName string) (ScriptFunc, error) {
	if f, ok := scriptRouteMap[scriptName]; ok {
		return f, nil
	} else {
		return nil, fmt.Errorf("script name %s not match", scriptName)
	}
}
