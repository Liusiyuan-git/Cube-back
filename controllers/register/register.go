package register

import (
	"Cube-back/controllers"
	"Cube-back/models/user"
	"Cube-back/models/user/profile"
	"Cube-back/models/user/register"
)

type Controller struct {
	controllers.MainController
}

var R = new(register.Register)
var P = new(profile.Profile)

func (c *Controller) UserRegister() {
	data := c.RequestBodyData()
	email := data["email"]
	password := data["password"]
	phone := data["phone"]
	code := data["code"]
	msg, pass := R.UserRegister(email, password, phone, code)
	result := make(map[string]interface{})
	result["msg"] = msg
	c.DataCallBack(result, pass)
}

func (c *Controller) VerificationCode() {
	data := c.RequestBodyData()
	phone := data["phone"]
	mode := data["mode"]
	u := new(user.User)
	content := u.VerificationCode(mode, phone)
	result := make(map[string]interface{})
	result["phone"] = phone
	result["content"] = content
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

func init() {
}
