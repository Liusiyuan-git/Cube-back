package profile

import (
	"Cube-back/controllers/login"
	"Cube-back/models/leaveMessage"
	"Cube-back/models/user/profile"
)

type Controller struct {
	login.Controller
}

var P = new(profile.Profile)
var l = new(leaveMessage.LeaveMessage)

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

func (c *Controller) UserNameSend() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	name := data["name"]
	pass := P.UserNameSend(cubeId, name)
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

func (c *Controller) UserImageUpdate() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	image, pass := P.UserImageUpdate(cubeId)
	result := make(map[string]interface{})
	result["image"] = image
	c.DataCallBack(result, pass)
}

func (c *Controller) UserCareSet() {
	data := c.RequestBodyData()
	id := data["id"]
	cubeId := data["cubeid"]
	msg, pass := P.UserCareSet(id, cubeId)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}

func (c *Controller) UserCareGet() {
	data := c.RequestBodyData()
	id := data["id"]
	cubeId := data["cubeid"]
	image, pass := P.UserCareGet(id, cubeId)
	result := make(map[string]interface{})
	result["image"] = image
	c.DataCallBack(result, pass)
}

func (c *Controller) UserCareConfirm() {
	data := c.RequestBodyData()
	id := data["id"]
	cubeId := data["cubeid"]
	exist, pass := P.UserCareConfirm(id, cubeId)
	result := make(map[string]interface{})
	result["exist"] = exist
	c.DataCallBack(result, pass)
}

func (c *Controller) UserCareCancel() {
	data := c.RequestBodyData()
	id := data["id"]
	cubeId := data["cubeid"]
	pass := P.UserCareCancel(id, cubeId)
	result := make(map[string]interface{})
	c.DataCallBack(result, pass)
}

func (c *Controller) ProfileLeaveSet() {
	data := c.RequestBodyData()
	cubeId := data["cubeId"]
	leaveId := data["leaveId"]
	text := data["text"]
	pass := l.LeaveSet(cubeId, leaveId, text)
	result := make(map[string]interface{})
	c.DataCallBack(result, pass)
}
