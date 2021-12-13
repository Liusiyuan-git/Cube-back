package talk

import (
	"Cube-back/controllers/login"
	"Cube-back/models/talk"
	"Cube-back/models/talkcomment"
)

type Controller struct {
	login.Controller
}

var b = new(talk.Talk)
var tc = new(talkcomment.TalkComment)

func (o *Controller) TalkSend() {
	data := o.RequestBodyData()
	text := data["text"]
	cubeid := data["cubeid"]
	images := data["images"]
	msg, pass := b.TalkSend(cubeid, text, images)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) TalkCommentSend() {
	data := o.RequestBodyData()
	text := data["text"]
	talkid := data["id"]
	cubeid := data["cubeid"]
	talkCubeId := data["talkCubeId"]
	msg, pass := tc.TalkCommentSend(talkid, cubeid, talkCubeId, text)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) TalkCommentDelete() {
	data := o.RequestBodyData()
	talkcommentid := data["id"]
	talkid := data["talkid"]
	cubeid := data["cubeid"]
	index := data["index"]
	msg, pass := tc.TalkCommentDelete(talkcommentid, cubeid, talkid, index)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) TalkDelete() {
	data := o.RequestBodyData()
	images := data["images"]
	talkId := data["id"]
	cubeId := data["cube_id"]
	index := data["index"]
	date := data["date"]
	msg, pass := b.TalkDelete(images, cubeId, talkId, index, date)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}
