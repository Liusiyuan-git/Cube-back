package community

import (
	"Cube-back/controllers/login"
)

type Controller struct {
	login.Controller
}

func (c *Controller) Test() {
	result := make(map[string]interface{})
	result["msg"] = "欢迎"
	c.DataCallBack(result, false)
}
