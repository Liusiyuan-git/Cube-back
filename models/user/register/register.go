package register

import (
	"Cube-back/database"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"github.com/go-basic/uuid"
)

type Register struct {
	user.User
}

func (r *Register) UserRegister(email, password, phone, code string) (string, bool) {
	s, pass := registerParamsCheck(email, password, phone, code)
	if !pass {
		return s, pass
	}
	cubeId := getCubeId()
	u := new(user.User)
	password = crypt.Set(password)
	u.Email, u.Password, u.Phone, u.CubeId = email, password, phone, cubeId
	_, err := database.Insert(u)
	if err != nil {
		return "注册失败，请稍后再试", false
	}
	return "", true
}

func registerParamsCheck(email, password, phone, code string) (string, bool) {
	var pass bool
	var msg string
	msg, pass = paramsEmpty(email, password, phone, code)
	if !pass {
		return msg, pass
	}
	msg, pass = paramsCorrect(email, phone, code)
	if !pass {
		return msg, pass
	}
	msg, pass = paramsRepeat(email, phone)
	if !pass {
		return msg, pass
	}
	return "", true
}

func paramsEmpty(email, password, phone, code string) (string, bool) {
	if email == "" || password == "" || phone == "" || code == "" {
		return "表单信息不完整，请检查", false
	}
	return "", true
}

func paramsCorrect(email, phone, code string) (string, bool) {
	msg, pPass := user.PhoneCorrect(phone)
	if !pPass {
		return msg, pPass
	}
	msg, ePass := user.EmailCorrect(email)
	if !ePass {
		return msg, ePass
	}
	msg, cPass := user.CodeCorrect(phone, code)
	if !cPass {
		return msg, cPass
	}
	return "", true
}

func paramsRepeat(email, phone string) (string, bool) {
	cmd := "select *  from user where email = ? or phone = ?"
	num, _, bool := database.DBValues(cmd, email, phone)
	if !bool {
		return "未知错误", false
	} else {
		if num > 0 {
			return "该 邮箱/手机号 已注册", false
		}
	}
	return "", true
}

func getCubeId() string {
	cubeId := uuid.New()
	return cubeId
}