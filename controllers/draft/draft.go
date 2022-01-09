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

func (o *Controller) DraftImageUpload() {
	data := o.RequestBodyData()
	cubeId := data["cube_id"]
	image := data["image"]
	mode := data["mode"]
	filename, message, pass := d.DraftImageUpload(cubeId, image, mode)
	result := make(map[string]interface{})
	result["filename"] = filename
	result["message"] = message
	o.DataCallBack(result, pass)
}

func (o *Controller) DraftImageDelete() {
	data := o.RequestBodyData()
	cubeId := data["cube_id"]
	filename := data["filename"]
	pass := d.DraftImageDelete(cubeId, filename)
	result := make(map[string]interface{})
	o.DataCallBack(result, pass)
}
