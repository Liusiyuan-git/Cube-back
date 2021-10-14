package blog

import (
	"Cube-back/controllers/login"
	"Cube-back/models/blog"
	"Cube-back/models/blogcomment"
	"Cube-back/models/collect"
)

type Controller struct {
	login.Controller
}

var b = new(blog.Blog)
var c = new(collect.Collect)
var bc = new(blogcomment.BlogComment)

func (o *Controller) BlogSend() {
	data := o.RequestBodyData()
	cover := data["cover"]
	title := data["title"]
	images := data["images"]
	content := data["content"]
	text := data["text"]
	cubeid := data["cubeid"]
	label := data["label"]
	labelType := data["labeltype"]
	msg, pass := b.BlogSend(cubeid, cover, title, content, text, images, label, labelType)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) BlogCollect() {
	data := o.RequestBodyData()
	cubeid := data["cubeid"]
	id := data["id"]
	cover := data["cover"]
	title := data["title"]
	date := data["date"]
	labelType := data["label_type"]
	msg, pass := c.BlogCollect(cubeid, id, cover, date, title, labelType)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) BlogCollectConfirm() {
	data := o.RequestBodyData()
	id := data["id"]
	cubeid := data["cubeid"]
	pass := c.BlogCollectConfirm(id, cubeid)
	result := make(map[string]interface{})
	o.DataCallBack(result, pass)
}

func (o *Controller) BlogCommentSend() {
	data := o.RequestBodyData()
	id := data["id"]
	cubeid := data["cubeid"]
	comment := data["comment"]
	msg, pass := bc.BlogCommentSend(id, cubeid, comment)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) BlogCollectionGet() {
	data := o.RequestBodyData()
	cubeid := data["cubeid"]
	content, length, pass := c.BlogCollectionGet(cubeid)
	result := make(map[string]interface{})
	result["content"] = content
	result["length"] = length
	o.DataCallBack(result, pass)
}
