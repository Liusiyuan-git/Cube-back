package draft

import (
	"Cube-back/controllers/login"
	"Cube-back/models/draft"
)

type Controller struct {
	login.Controller
}

var d = new(draft.Draft)

func (o *Controller) DraftSend() {
	data := o.RequestBodyData()
	cover := data["cover"]
	title := data["title"]
	images := data["images"]
	content := data["content"]
	cubeid := data["cubeid"]
	msg, pass := d.DraftSend(cubeid, cover, title, content, images)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) DraftGet() {
	data := o.RequestBodyData()
	cubeid := data["cubeid"]
	content, pass := d.DraftGet(cubeid)
	result := make(map[string]interface{})
	result["content"] = content
	o.DataCallBack(result, pass)
}

func (o *Controller) DraftRemove() {
	data := o.RequestBodyData()
	cubeid := data["cubeid"]
	content, pass := d.DraftRemove(cubeid)
	result := make(map[string]interface{})
	result["content"] = content
	o.DataCallBack(result, pass)
}
