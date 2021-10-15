package blogcomment

import (
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"strconv"
)

func BlogCommentRedisGet(id, page string) (interface{}, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var dataBlock []map[string]interface{}
	var d = redis.LRange("blog_detail_comment_"+id, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen("blog_detail_comment_" + id)
	for _, item := range d {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	redis.HSet("blog_profile_"+id, "comment", fmt.Sprintf("%v", l))
	return dataBlock, l
}

func blogCommentRedisLock(key, status string) {
	redis.Set(key, status)
}

func blogCommentRedisLockStatus(key string) string {
	status := redis.Get(key)
	return status
}

func blogCommonLikeRedis(commentid, blogid, index, love string) {
	i, _ := strconv.Atoi(index)
	key := "blog_detail_comment_" + blogid
	comment := redis.LIndex(key, int64(i))
	if comment != "" {
		var m map[string]interface{}
		json.Unmarshal([]byte(comment), &m)
		if m["id"] == commentid {
			m["love"] = love
			bjson, _ := json.Marshal(m)
			redisValue := string(bjson)
			redis.LSet(key, int64(i), redisValue)
		}
	}
}

func blogCommentSendDbRedis(blogid, cubeid, comment, date string, commentId int64) {
	b := make(map[string]interface{})
	key := "blog_detail_comment_" + blogid
	userName := redis.HGet("session", cubeid)
	b["id"] = strconv.FormatInt(commentId, 10)
	b["cube_id"] = cubeid
	b["comment"] = comment
	b["date"] = date
	b["love"] = "0"
	b["name"] = userName
	b["image"] = redis.HGet("user_profile_"+cubeid, "image")
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush(key, redisValue)
	redis.HIncrBy("blog_profile_"+blogid, "comment", 1)
}
