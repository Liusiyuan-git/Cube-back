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
	msg, pass := b.TalkSend(cubeid, text)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) TalkCommentSend() {
	data := o.RequestBodyData()
	text := data["text"]
	talkid := data["id"]
	cubeid := data["cubeid"]
	count := data["comment"]
	msg, pass := tc.TalkCommentSend(talkid, cubeid, text, count)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}

func (o *Controller) TalkCommentDelete() {
	data := o.RequestBodyData()
	talkcommentid := data["id"]
	talkid := data["talkid"]
	cubeid := data["cubeid"]
	count := data["comment"]
	msg, pass := tc.TalkCommentDelete(talkcommentid, cubeid, talkid, count)
	result := make(map[string]interface{})
	result["msg"] = msg
	o.DataCallBack(result, pass)
}
