package blogcomment

import (
	"Cube-back/rabbitmq"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	Redis "github.com/go-redis/redis"
	"math"
	"strconv"
)

func BlogCommentRedisGet(id, page string) (interface{}, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var dataBlock []map[string]interface{}
	var l = redis.LLen("blog_detail_comment_" + id)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var d = redis.LRange("blog_detail_comment_"+id, (pageInt-1)*10, (pageInt-1)*10+9)
	for _, item := range d {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	redis.HSet("blog_profile_"+id, "comment", fmt.Sprintf("%v", l))
	return dataBlock, l
}

func blogCommentRedisLock(key, status string) {
	redis.Lock(key, status)
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
	txpipeline := redis.TxPipeline()
	txpipeline.LPush(key, redisValue)
	txpipeline.HIncrBy("blog_profile_"+blogid, "comment", 1)
	txpipeline.Exec()
	txpipeline.Close()
}

func blogCommentDeleteRedisUpdate(blogCommentId, blogId, index string) {
	var key = "blog_detail_comment_" + blogId
	var m map[string]interface{}
	location, _ := strconv.Atoi(index)
	each := redis.LIndex(key, int64(location))
	json.Unmarshal([]byte(each), &m)
	if blogCommentId == m["id"] {
		redis.LRem(key, each)
	} else {
		blogCommentBox := redis.LRange(key, 0, -1)
		for _, item := range blogCommentBox {
			var s map[string]interface{}
			json.Unmarshal([]byte(item), &s)
			if blogCommentId == s["id"] {
				redis.LRem(key, item)
				break
			}
		}
	}
}

func blogCommentMessageSendRedis(blogCubeId, cubeid string, messageId int64, bc *BlogComment) {
	b := make(map[string]interface{})
	b["send_id"] = cubeid
	b["text"] = bc.Comment
	b["id"] = strconv.FormatInt(messageId, 10)
	b["date"] = bc.Date
	b["blog_comment"] = "1"
	b["blog_id"] = strconv.Itoa(bc.BlogId)
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	txpipeline := redis.TxPipeline()
	execBlock := []interface{}{
		txpipeline.HIncrBy("user_message_profile_"+blogCubeId, "total", 1),
		txpipeline.LPush("user_message_"+blogCubeId, redisValue),
	}
	txpipeline.Exec()
	rabbitmq.MessageQueue.MessageSend(blogCubeId, fmt.Sprintf("%v", execBlock[0].(*Redis.IntCmd).Val()))
	txpipeline.Close()
}
