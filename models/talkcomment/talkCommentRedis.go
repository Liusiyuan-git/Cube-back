package talkcomment

import (
	"Cube-back/redis"
	"encoding/json"
	"strconv"
)

func TalkCommentRedisGet(talkId, page string) (interface{}, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var dataBlock []map[string]interface{}
	var d = redis.LRange("talk_comment_"+talkId, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen("talk_comment_" + talkId)
	for _, item := range d {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, l
}

func TalkCommentRedisSend(id int64, talkid, index, mode, commentCount string, talkComment TalkComment) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", talkComment.CubeId)
	b["id"] = id
	b["name"] = userName
	b["talk_id"] = talkid
	b["cube_id"] = talkComment.CubeId
	b["comment"] = talkComment.Comment
	b["date"] = talkComment.Date
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush("talk_comment_"+talkid, redisValue)
	talkCommentNumRedis(talkid, index, mode, commentCount)
}

func talkCommentNumRedis(talkid, index, mode, commentCount string) {
	i, _ := strconv.Atoi(index)
	key := "talk_" + mode
	comment := redis.LIndex(key, int64(i))
	if comment != "" {
		var m map[string]interface{}
		json.Unmarshal([]byte(comment), &m)
		if m["id"] == talkid {
			m["comment"] = commentCount
			bjson, _ := json.Marshal(m)
			redisValue := string(bjson)
			redis.LSet(key, int64(i), redisValue)
		}
	}
}

func talkCommentRedisLockStatus(key string) string {
	status := redis.Get(key)
	return status
}

func talkCommentRedisLock(key, status string) {
	redis.Set(key, status)
}

func talkDetailRedisSet(talkId string) {
	redis.HSet("talk_detail", talkId, "1")
}
