package login

import (
	"Cube-back/controllers"
	"Cube-back/log"
	"Cube-back/models/user/login"
)

type Controller struct {
	controllers.MainController
}

func (c *Controller) Prepare() {
	data := c.RequestBodyData()
	mode := data["mode"]
	if mode == "" {
		session := c.GetSession("CubeId")
		if session == nil {
			result := make(map[string]interface{})
			result["msg"] = "请先登录"
			c.DataCallBack(result, false)
		}
	}
}

var L = new(login.Login)

func (c *Controller) UserLogin() {
	data := c.RequestBodyData()
	mode := data["mode"]
	cubeId := ""
	msg := ""
	pass := true
	if mode == "count" {
		count := data["count"]
		password := data["password"]
		cubeId, msg, pass = L.LoginCount(count, password)
	}
	if mode == "phone" {
		phone := data["phone"]
		code := data["code"]
		cubeId, msg, pass = L.LoginPhone(phone, code)
	}
	if pass {
		sessionErr := c.SetSession("CubeId", cubeId)
		if sessionErr != nil {
			log.Error(sessionErr)
		}
	}
	result := make(map[string]interface{})
	result["cubeId"] = cubeId
	result["msg"] = msg
	c.DataCallBack(result, pass)
}
