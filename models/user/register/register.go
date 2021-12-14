package register

import (
	"Cube-back/database"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"Cube-back/rabbitmq"
	"Cube-back/redis"
	"Cube-back/snowflake"
	"fmt"
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
	u.Email, u.Password, u.Phone, u.CubeId, u.Name = email, password, phone, cubeId, phone[0:3]+"****"+phone[7:]
	_, err := database.Insert(u)
	if err != nil {
		return "注册失败，请稍后再试", false
	}
	userRedisProfile(u.CubeId, u.Phone)
	rabbitmq.MessageQueue.MessageSend(cubeId, fmt.Sprintf("%v", "欢迎来到cube"))
	return "", true
}

func userRedisProfile(cubeId, phone string) {
	key := "user_profile_" + cubeId
	txpipeline := redis.TxPipeline()
	txpipeline.HSet(key, "name", phone[0:3]+"****"+phone[7:])
	txpipeline.HIncrBy(key, "blog", 0)
	txpipeline.HIncrBy(key, "blog", 0)
	txpipeline.HIncrBy(key, "talk", 0)
	txpipeline.HIncrBy(key, "collect", 0)
	txpipeline.HIncrBy(key, "cared", 0)
	txpipeline.HIncrBy(key, "care", 0)
	txpipeline.HIncrBy(key, "leave", 0)
	txpipeline.HIncrBy(key, "message", 0)
	txpipeline.Exec()
	txpipeline.Close()
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
	cubeId := snowflake.NewNode.NewId()
	return cubeId
}
