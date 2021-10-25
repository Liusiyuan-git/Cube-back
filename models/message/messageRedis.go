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
	return redis.HMGet("user_message_profile_"+cubeId, []string{"total", "blog", "talk"})
}

func UserMessageCleanRedis(id string) {
	var result string
	box := redis.HMGet("user_message_profile_"+id, []string{"blog", "talk"})
	blog, _ := strconv.Atoi(box[0].(string))
	talk, _ := strconv.Atoi(box[1].(string))
	value := blog + talk
	result = strconv.Itoa(value)
	redis.HSet("user_message_profile_"+id, "total", result)
}

func messageProfileUserTalkRedisGet(cubeId string, ids []string) interface{} {
	return redis.HMGet("user_message_profile_"+cubeId, ids)
}

func messageProfileUserTalkRedisClean(id, deleteId string) {
	key := "user_message_profile_" + id
	target := "talk_" + deleteId
	userTalk := redis.HGet(key, target)
	if userTalk != "nil" {
		t, _ := strconv.Atoi(userTalk)
		r := int64(-1 * t)
		redis.HIncrBy(key, "talk", r)
		redis.HIncrBy(key, "total", r)
		redis.HDel(key, target)
	}
}
