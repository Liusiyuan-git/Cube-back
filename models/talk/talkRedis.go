package talk

import (
	"Cube-back/redis"
	"encoding/json"
	"strconv"
)

func talkRedisGet(mode, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var key = "talk_" + mode
	var t = redis.LRange(key, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen(key)
	return t, l
}

func talkSendRedis(id int64, talk Talk) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", talk.CubeId)
	b["id"] = id
	b["name"] = userName
	b["cube_id"] = talk.CubeId
	b["comment"] = talk.Comment
	b["text"] = talk.Text
	b["love"] = talk.Love
	b["date"] = talk.Date
	b["images"] = talk.Images
	b["user_image"] = redis.HGet("user_profile_"+talk.CubeId, "image")
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush("talk_new", redisValue)
	redis.LPush("profile_talk_"+talk.CubeId, redisValue)
	redis.HIncrBy("user_profile_"+talk.CubeId, "talk", 1)
	redis.HIncrBy("talk_like_and_comment", strconv.FormatInt(id, 10)+"_like", 0)
	redis.HIncrBy("talk_like_and_comment", strconv.FormatInt(id, 10)+"_comment", 0)
}

func TalkRedisLockStatus(key string) string {
	status := redis.Get(key)
	return status
}

func TalkLikeRedis(talkid, like, index, mode string) {
	i, _ := strconv.Atoi(index)
	key := "talk_" + mode
	comment := redis.LIndex(key, int64(i))
	if comment != "" {
		var m map[string]interface{}
		json.Unmarshal([]byte(comment), &m)
		if m["id"] == talkid {
			m["love"] = like
			bjson, _ := json.Marshal(m)
			redisValue := string(bjson)
			redis.LSet(key, int64(i), redisValue)
		}
	}
}

func TalkLikeRedisProfile(talkid, cubeId, like, index string) {
	i, _ := strconv.Atoi(index)
	key := "profile_talk_" + cubeId
	comment := redis.LIndex(key, int64(i))
	if comment != "" {
		var m map[string]interface{}
		json.Unmarshal([]byte(comment), &m)
		if m["id"] == talkid {
			m["love"] = like
			bjson, _ := json.Marshal(m)
			redisValue := string(bjson)
			redis.LSet(key, int64(i), redisValue)
		}
	}
}

func TalkRedisLock(key, status string) {
	redis.Set(key, status)
}
