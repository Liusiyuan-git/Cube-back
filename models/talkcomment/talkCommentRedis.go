package talkcomment

import (
	"Cube-back/rabbitmq"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

func TalkCommentRedisGet(talkId, page string) (interface{}, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var dataBlock []map[string]interface{}
	var l = redis.LLen("talk_comment_" + talkId)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var d = redis.LRange("talk_comment_"+talkId, (pageInt-1)*10, (pageInt-1)*10+9)
	for _, item := range d {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, l
}

func TalkCommentRedisSend(id int64, talkid string, talkComment TalkComment) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", talkComment.CubeId)
	b["id"] = strconv.FormatInt(id, 10)
	b["name"] = userName
	b["talk_id"] = talkid
	b["cube_id"] = talkComment.CubeId
	b["comment"] = talkComment.Comment
	b["date"] = talkComment.Date
	b["user_image"] = redis.HGet("user_profile_"+talkComment.CubeId, "image")
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush("talk_comment_"+talkid, redisValue)

}

func talkCommentRedisLockStatus(key string) string {
	status := redis.Get(key)
	return status
}

func talkCommentRedisLock(key, status string) {
	redis.Set(key, status)
}

func talkCommentDeleteRedisUpdate(talkcommentid, talkId, index string) {
	var key = "talk_comment_" + talkId
	var m map[string]interface{}
	location, _ := strconv.Atoi(index)
	each := redis.LIndex(key, int64(location))
	json.Unmarshal([]byte(each), &m)
	if talkcommentid == m["id"] {
		redis.LRem(key, each)
		redis.HIncrBy("talk_like_and_comment", talkId+"_comment", -1)
	} else {
		talkCommentBox := redis.LRange(key, 0, -1)
		for _, item := range talkCommentBox {
			var s map[string]interface{}
			json.Unmarshal([]byte(item), &s)
			if talkcommentid == s["id"] {
				redis.LRem(key, item)
				redis.HIncrBy("talk_like_and_comment", talkId+"_comment", -1)
				break
			}
		}
	}
}

func userCareRedisGet(cubeId string) []string {
	return redis.HKeys("user_cared_" + cubeId)
}

func talkCommentMessageSendRedis(talkCubeId, cubeid string, messageId int64, tc *TalkComment) {
	b := make(map[string]interface{})
	b["send_id"] = cubeid
	b["text"] = tc.Comment
	b["id"] = strconv.FormatInt(messageId, 10)
	b["date"] = tc.Date
	b["talk_comment"] = "1"
	b["talk_id"] = strconv.Itoa(tc.TalkId)
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush("user_message_"+talkCubeId, redisValue)
	rabbitmq.MessageQueue.MessageSend(talkCubeId, fmt.Sprintf("%v", redis.HIncrBy("user_message_profile_"+talkCubeId, "total", 1)))
}
