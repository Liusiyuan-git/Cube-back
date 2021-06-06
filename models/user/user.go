package user

import (
	"Cube-back/database"
	"Cube-back/redis"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type User struct {
	Id       int
	CubeId   string
	Name     string
	Email    string
	Password string
	Phone    string
}

func random() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}

func PhoneCorrect(phone string) (string, bool) {
	regular := "^1[0-9]{10}$"
	reg := regexp.MustCompile(regular)
	correct := reg.MatchString(phone)
	if correct {
		return "", true
	}
	return "手机号格式异常", false
}

func EmailCorrect(email string) (string, bool) {
	pattern := ""
	reg := regexp.MustCompile(pattern)
	correct := reg.MatchString(email)
	if correct {
		return "", true
	}
	return "邮箱格式异常", false
}

func CodeCorrect(phone, code string) (string, bool) {
	store := redis.HGet("VerificationCode", phone)
	if code == store {
		return "", true
	}
	return "验证码错误", false
}

func (u *User) VerificationCode(phone string) {
	value := random()
	redis.HSet("VerificationCode", phone, value)
}

func PhoneConfirm(phone string) (string, bool) {
	cmd := "select * from user where phone = ?"
	num, maps, pass := database.DBValues(cmd, phone)
	if !pass {
		return "未知错误", false
	} else {
		if num > 0 && maps[0]["phone"] == phone {
			id := fmt.Sprintf("%v", maps[0]["id"])
			return id, true
		} else {
			return "手机号为空或不存在", false
		}
	}
}
