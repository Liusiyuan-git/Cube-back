package user

import (
	"Cube-back/controllers"
	"Cube-back/models/collect"
	"Cube-back/models/user"
	"Cube-back/models/user/profile"
	"fmt"
	"io/ioutil"
	"path"
)

type Controller struct {
	controllers.MainController
}

var P = new(profile.Profile)
var collection = new(collect.Collect)

func (c *Controller) VerificationCode() {
	data := c.RequestBodyData()
	phone := data["phone"]
	u := new(user.User)
	u.VerificationCode(phone)
	result := make(map[string]interface{})
	result["phone"] = phone
	c.DataCallBack(result, true)
}

func (c *Controller) PasswordChange() {
	data := c.RequestBodyData()
	phone := data["phone"]
	password := data["password"]
	code := data["code"]
	msg, pass := P.PasswordChange(phone, password, code)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}

func (c *Controller) GetUserPhoto() {
	filename := c.GetString("userid")
	kind := c.GetString("kind")
	typeFile := c.GetString("type")
	imgBase := path.Join("picture/user", "user.jpg")
	img := path.Join("picture/"+kind, filename+"."+typeFile)

	c.Ctx.Output.Header("Content-Type", "image/jpg")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", img))
	file, err := ioutil.ReadFile(img)
	if err != nil {
		file, _ := ioutil.ReadFile(imgBase)
		c.Ctx.WriteString(string(file))
	} else {
		c.Ctx.WriteString(string(file))
	}
}

func init() {
}
