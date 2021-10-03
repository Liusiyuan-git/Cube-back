package profile

import (
	"Cube-back/controllers/login"
	"Cube-back/models/user/profile"
)

type Controller struct {
	login.Controller
}

var P = new(profile.Profile)

func (c *Controller) SendUserImage() {
	data := c.RequestBodyData()
	image := data["image"]
	cubeId := data["cubeid"]
	msg, pass := P.SendUserImage(cubeId, image)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}

func (c *Controller) UserIntroduceSend() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	introduce := data["introduce"]
	pass := P.UserIntroduceSend(cubeId, introduce)
	result := make(map[string]interface{})
	c.DataCallBack(result, pass)
}

func (c *Controller) UserProfileGet() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	profile, pass := P.UserProfileGet(cubeId)
	result := make(map[string]interface{})
	result["profile"] = profile
	c.DataCallBack(result, pass)
}
