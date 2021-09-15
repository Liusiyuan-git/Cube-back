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

func (c *Controller) LoginStatus() {
	data := c.RequestBodyData()
	mode := data["mode"]
	if mode == "" {
		session := c.GetSession("CubeId")
		status := true
		if session == nil {
			status = false
		}
		result := make(map[string]interface{})
		result["cubeId"] = session
		c.DataCallBack(result, status)
	}
}

func (c *Controller) CountExit() {
	err := c.DelSession("CubeId")
	result := make(map[string]interface{})
	if err != nil {
		log.Error(err)
	}
	c.DataCallBack(result, err == nil)
}

var L = new(login.Login)

func (c *Controller) UserLogin() {
	data := c.RequestBodyData()
	mode := data["mode"]
	cubeId := ""
	userName := ""
	image := ""
	msg := ""
	pass := true
	if mode == "count" {
		count := data["count"]
		password := data["password"]
		cubeId, msg, userName, image, pass = L.LoginCount(count, password)
	}
	if mode == "phone" {
		phone := data["phone"]
		code := data["code"]
		cubeId, userName, msg, image, pass = L.LoginPhone(phone, code)
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
	result["userName"] = userName
	result["image"] = image
	c.DataCallBack(result, pass)
}
