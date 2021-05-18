package log

import (
	"github.com/beego/beego/v2/core/logs"
)

func Error(e interface{}) {
	logs.Error(e)
}

func Info(s interface{}) {
	logs.Info(s)
}

func Warn(s interface{}) {
	logs.Warn(s)
}

func init() {
	_ = logs.SetLogger(logs.AdapterFile, `{"filename":"log/cube.log","color":true, "level":6}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
	logs.Info("info")
}
