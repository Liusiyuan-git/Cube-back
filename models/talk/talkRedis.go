package talk

import (
	"Cube-back/rabbitmq"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
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

func userCareRedisGet(cubeId string) []string {
	return redis.HKeys("user_cared_" + cubeId)
}

func talkMessageSendRedis(t *Talk) {
	var userTalkText string
	b := make(map[string]interface{})
	caredBox := userCareRedisGet(t.CubeId)
	for _, item := range caredBox {
		userName := redis.HGet("session", item)
		userImage := redis.HGet("user_profile_"+item, "image")
		if len(t.Text) > 30 {
			userTalkText = t.Text[30:]
		} else {
			userTalkText = t.Text
		}
		b["cube_id"] = item
		b["send_id"] = t.CubeId
		b["date"] = t.Date
		b["text"] = userTalkText
		b["talk"] = 1
		b["name"] = userName
		b["image"] = userImage
		bjson, _ := json.Marshal(b)
		redisValue := string(bjson)
		redis.LPush("user_message_"+item, redisValue)
		redis.HIncrBy("user_message_profile_"+item, "talk", 1)
		redis.HIncrBy("user_message_profile_"+item, "talk_"+t.CubeId, 1)
		rabbitmq.MessageQueue.MessageSend(item, fmt.Sprintf("%v", redis.HIncrBy("user_message_profile_"+item, "total", 1)))
	}
}
