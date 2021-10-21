package login

import (
	"Cube-back/database"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"Cube-back/redis"
	"fmt"
)

type Login struct {
	user.User
}

func (u *Login) LoginCount(count, password string) (string, string, string, string, bool) {
	cubeId, p, userName, image, pass := u.CountConfirm(count)
	if !pass {
		return cubeId, p, "", "", false
	}
	pass = crypt.Confirm(password, p)
	if !pass {
		return cubeId, "密码错误", "", "", false
	}
	return cubeId, "", userName, image, true
}

func (u *Login) LoginPhone(phone, code string) (string, string, string, string, bool) {
	msg, pass := user.CodeCorrect(phone, code)
	if !pass {
		return "", "", msg, "", false
	}
	cubeId, userName, msg, image, pass := PhoneConfirm(phone)
	if !pass {
		return "", "", msg, "", false
	}
	return cubeId, userName, "", image, true
}

func (u *Login) CountConfirm(count string) (string, string, string, string, bool) {
	cmd := "select * from user where email = ?"
	num, maps, pass := database.DBValues(cmd, count)
	if !pass {
		return "", "未知错误", "", "", false
	} else {
		if num > 0 && (maps[0]["email"] == count || maps[0]["name"] == count) {
			password := fmt.Sprintf("%v", maps[0]["password"])
			cubeId := fmt.Sprintf("%v", maps[0]["cube_id"])
			userName := fmt.Sprintf("%v", maps[0]["name"])
			image := fmt.Sprintf("%v", maps[0]["image"])
			sessionRedis(cubeId, userName)
			userImageRedis(cubeId, image)
			userMessageRedis(cubeId)
			return cubeId, password, userName, image, true
		} else {
			return "", "账号不存在", "", "", false
		}
	}
}

func userMessageRedis(cubeId string) {
	redis.HIncrBy("user_message_profile_"+cubeId, "total", 1)
	redis.HIncrBy("user_message_profile_"+cubeId, "blog", 0)
	redis.HIncrBy("user_message_profile_"+cubeId, "talk", 0)
}

func sessionRedis(cubeid, name string) {
	redis.HSet("session", cubeid, name)
}

func userImageRedis(cubeId, image string) {
	redis.HSet("userImage", cubeId, image)
}

func PhoneConfirm(phone string) (string, string, string, string, bool) {
	cmd := "select * from user where phone = ?"
	num, maps, pass := database.DBValues(cmd, phone)
	if !pass {
		return "", "", "未知错误", "", false
	} else {
		if num > 0 && maps[0]["phone"] == phone {
			cubeId := fmt.Sprintf("%v", maps[0]["cube_id"])
			userName := fmt.Sprintf("%v", maps[0]["name"])
			image := fmt.Sprintf("%v", maps[0]["image"])
			return cubeId, userName, "", image, true
		} else {
			return "", "", "手机号不存在", "", false
		}
	}
}
