package private

import (
	"Cube-back/controllers/login"
	"Cube-back/models/blog"
	"Cube-back/models/collect"
	"Cube-back/models/draft"
)

type Controller struct {
	login.Controller
}

var b = new(blog.Blog)
var d = new(draft.Draft)
var c = new(collect.Collect)

func (o *Controller) Test() {
	result := make(map[string]interface{})
	result["msg"] = "欢迎"
	o.DataCallBack(result, false)
}
