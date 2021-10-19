package message

import (
	"Cube-back/redis"
	"strconv"
)

func userMessageRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var key = "user_message_" + cubeId
	var t = redis.LRange(key, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen(key)
	return t, l
}

func MessageProfileRedisGet(cubeId string) interface{} {
	return redis.HMGet("user_message_profile_"+cubeId, []string{"total"})
}
