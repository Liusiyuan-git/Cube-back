package login

import (
	"Cube-back/database"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"fmt"
)

type Login struct {
	user.User
}

func (u *Login) LoginCount(count, password string) (string, string, bool) {
	cubeId, p, pass := u.CountConfirm(count)
	if !pass {
		return cubeId, p, false
	}
	pass = crypt.Confirm(password, p)
	if !pass {
		return cubeId, "密码错误", false
	}
	return cubeId, "", true
}

func (u *Login) LoginPhone(phone, code string) (string, string, bool) {
	msg, pass := user.CodeCorrect(phone, code)
	if !pass {
		return "", msg, false
	}
	cubeId, msg, pass := PhoneConfirm(phone)
	if !pass {
		return "", msg, false
	}
	return cubeId, "", true
}

func (u *Login) CountConfirm(count string) (string, string, bool) {
	cmd := "select * from user where email = ? or name = ?"
	num, maps, pass := database.DBValues(cmd, count, count)
	if !pass {
		return "", "未知错误", false
	} else {
		if num > 0 && (maps[0]["email"] == count || maps[0]["name"] == count) {
			password := fmt.Sprintf("%v", maps[0]["password"])
			cubeId := fmt.Sprintf("%v", maps[0]["cube_id"])
			return cubeId, password, true
		} else {
			return "", "账号不存在", false
		}
	}
}

func PhoneConfirm(phone string) (string, string, bool) {
	cmd := "select * from user where phone = ?"
	num, maps, pass := database.DBValues(cmd, phone)
	if !pass {
		return "", "未知错误", false
	} else {
		if num > 0 && maps[0]["phone"] == phone {
			cubeId := fmt.Sprintf("%v", maps[0]["cube_id"])
			return cubeId, "", true
		} else {
			return "", "手机号不存在", false
		}
	}
}
