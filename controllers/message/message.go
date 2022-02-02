package message

import (
	"Cube-back/controllers/login"
	"Cube-back/models/message"
)

type Controller struct {
	login.Controller
}

var m = new(message.Message)

func (c *Controller) UserMessageGet() {
	data := c.RequestBodyData()
	cubeId := data["id"]
	page := data["page"]
	profile, length, pass := m.UserMessageGet(cubeId, page)
	result := make(map[string]interface{})
	result["content"] = profile
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) MessageProfileGet() {
	data := c.RequestBodyData()
	cubeId := data["cube_id"]
	profile := m.MessageProfileGet(cubeId)
	result := make(map[string]interface{})
	result["profile"] = profile
	c.DataCallBack(result, true)
}

func (c *Controller) UserMessageClean() {
	data := c.RequestBodyData()
	cubeId := data["id"]
	m.UserMessageClean(cubeId)
	result := make(map[string]interface{})
	c.DataCallBack(result, true)
}

func (c *Controller) MessageProfileUserTalkGet() {
	data := c.RequestBodyData()
	id := data["id"]
	idBox := data["idBox"]
	content, pass := m.MessageProfileUserTalkGet(id, idBox)
	result := make(map[string]interface{})
	result["content"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) MessageProfileUserTalkClean() {
	data := c.RequestBodyData()
	id := data["id"]
	deleteId := data["deleteId"]
	m.MessageProfileUserTalkClean(id, deleteId)
	result := make(map[string]interface{})
	c.DataCallBack(result, true)
}

func (c *Controller) MessageProfileUserBlogClean() {
	data := c.RequestBodyData()
	id := data["id"]
	deleteId := data["deleteId"]
	m.MessageProfileUserBlogClean(id, deleteId)
	result := make(map[string]interface{})
	c.DataCallBack(result, true)
}

func (c *Controller) MessageProfileUserBlogGet() {
	data := c.RequestBodyData()
	id := data["id"]
	idBox := data["idBox"]
	content, pass := m.MessageProfileUserBlogGet(id, idBox)
	result := make(map[string]interface{})
	result["content"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) MessageDelete() {
	data := c.RequestBodyData()
	id := data["id"]
	cubeId := data["cube_id"]
	index := data["index"]
	msg, pass := m.MessageDelete(id, cubeId, index)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}
