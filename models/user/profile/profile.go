package profile

import (
	"Cube-back/database"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"fmt"
	"strconv"
)

type Profile struct {
	user.User
}

func (p *Profile) PasswordChange(phone, password, code string) (string, bool) {
	s, pass, id := passwordParamsCheck(phone, password, code)
	if !pass {
		return s, pass
	}
	idi, err := strconv.Atoi(id)
	password = crypt.Set(password)
	u := new(user.User)
	u.Id, u.Password = idi, password
	_, err = database.Update(u, "password")
	if err != nil {
		fmt.Println(err)
		return "密码修改失败", false
	}
	return "", true
}

func passwordParamsCheck(phone, password, code string) (string, bool, string) {
	var pass bool
	var msg string
	msg, pass = paramsEmpty(password, phone, code)
	if !pass {
		return msg, pass, ""
	}
	msg, pass = user.CodeCorrect(phone, code)
	if !pass {
		return msg, pass, ""
	}
	msg, pass = user.PhoneConfirm(phone)
	if !pass {
		return msg, pass, ""
	}
	return "", true, msg
}

func paramsEmpty(password, phone, code string) (string, bool) {
	if password == "" || phone == "" || code == "" {
		return "表单信息不完整，请检查", false
	}
	return "", true
}
