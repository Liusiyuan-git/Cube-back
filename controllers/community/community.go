package community

import (
	"Cube-back/controllers/login"
	"Cube-back/models/blog"
	"Cube-back/models/draft"
)

type Controller struct {
	login.Controller
}

var b = new(blog.Blog)
var d = new(draft.Draft)

func (c *Controller) Test() {
	result := make(map[string]interface{})
	result["msg"] = "欢迎"
	c.DataCallBack(result, false)
}

func (c *Controller) BlogSend() {
	data := c.RequestBodyData()
	cover := data["cover"]
	title := data["title"]
	images := data["images"]
	content := data["content"]
	text := data["text"]
	cubeid := data["cubeid"]
	msg, pass := b.BlogSend(cubeid, cover, title, content, text, images)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}

func (c *Controller) DraftSend() {
	data := c.RequestBodyData()
	cover := data["cover"]
	title := data["title"]
	images := data["images"]
	content := data["content"]
	cubeid := data["cubeid"]
	msg, pass := d.DraftSend(cubeid, cover, title, content, images)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}

func (c *Controller) DraftGet() {
	data := c.RequestBodyData()
	cubeid := data["cubeid"]
	content, pass := d.DraftGet(cubeid)
	result := make(map[string]interface{})
	result["content"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) DraftRemove() {
	data := c.RequestBodyData()
	cubeid := data["cubeid"]
	content, pass := d.DraftRemove(cubeid)
	result := make(map[string]interface{})
	result["content"] = content
	c.DataCallBack(result, pass)
}
