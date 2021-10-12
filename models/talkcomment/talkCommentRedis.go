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

func TalkCommentRedisSend(id int64, talkid string, talkComment TalkComment) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", talkComment.CubeId)
	b["id"] = id
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

func talkDetailRedisSet(talkId string) {
	redis.HSet("talk_detail", talkId, "1")
}

func talkCommentDeleteRedisUpdate(talkId, index string) {
	var key = "talk_comment_" + talkId
	location, _ := strconv.Atoi(index)
	d := redis.LIndex(key, int64(location))
	redis.LRem(key, d)
}
