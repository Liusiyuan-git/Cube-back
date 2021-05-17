package controllers

import (
	"encoding/json"
	//"github.com/beego/beego/v2/server/web/context"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) RequestBodyData() (d map[string]string) {
	var data map[string]string
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &data)
	return data
}

func (c *MainController) DataCallBack(params map[string]interface{}, pass bool) {
	if !pass {
		params["success"] = false
	} else {
		params["success"] = true
	}
	c.Data["json"] = params
	_ = c.ServeJSON()
}
