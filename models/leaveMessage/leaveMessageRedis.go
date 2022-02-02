package leaveMessage

import (
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
	"math"
	"strconv"
)

func leaveRedisSet(cubeId, leaveId string, insertId int64, l *LeaveMessage) {
	profile := redis.HMGet("user_profile_"+leaveId, []string{"image", "name"})
	var leave = map[string]string{}
	leave["id"] = strconv.FormatInt(insertId, 10)
	leave["text"] = l.Text
	leave["leave_id"] = l.LeaveId
	leave["image"] = profile[0].(string)
	leave["name"] = profile[1].(string)
	leave["date"] = l.Date
	bjson, _ := json.Marshal(leave)
	redisValue := string(bjson)
	redis.LPush("user_leave_"+cubeId, redisValue)
}

func leaveRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var key = "user_leave_" + cubeId
	var l = redis.LLen(key)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var t = redis.LRange(key, (pageInt-1)*10, (pageInt-1)*10+9)
	return t, l
}

func leaveMessageRedisSet(insertId int64, message *message.Message) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", message.SendId)
	userImage := redis.HGet("user_profile_"+message.SendId, "image")
	b["id"] = strconv.FormatInt(insertId, 10)
	b["cube_id"] = message.CubeId
	b["send_id"] = message.SendId
	b["date"] = message.Date
	b["text"] = message.Text
	b["care"] = message.Care
	b["name"] = userName
	b["image"] = userImage
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush("user_message_"+message.CubeId, redisValue)
}

func leaveDeleteRedis(id, index, cubeId string) {
	var key = "user_leave_" + cubeId
	var m map[string]interface{}
	location, _ := strconv.Atoi(index)
	each := redis.LIndex(key, int64(location))
	json.Unmarshal([]byte(each), &m)
	if id == m["id"] {
		redis.LRem(key, each)
	} else {
		leaveBox := redis.LRange(key, 0, -1)
		for _, item := range leaveBox {
			var s map[string]interface{}
			json.Unmarshal([]byte(item), &s)
			if id == s["id"] {
				redis.LRem(key, item)
				break
			}
		}
	}
}
