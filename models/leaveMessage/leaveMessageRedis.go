package leaveMessage

import (
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
	"strconv"
)

func leaveRedisSet(cubeId, leaveId string, l *LeaveMessage) {
	profile := redis.HMGet("user_profile_"+leaveId, []string{"image", "name"})
	var leave = map[string]interface{}{}
	leave["text"] = l.Text
	leave["leave_id"] = l.LeaveId
	leave["image"] = profile[0]
	leave["name"] = profile[1]
	leave["date"] = l.Date
	bjson, _ := json.Marshal(leave)
	redisValue := string(bjson)
	redis.LPush("user_leave_"+cubeId, redisValue)
}

func leaveRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var key = "user_leave_" + cubeId
	var t = redis.LRange(key, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen(key)
	return t, l
}

func leaveMessageRedisSet(message *message.Message) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", message.SendId)
	userImage := redis.HGet("user_profile_"+message.SendId, "image")
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
