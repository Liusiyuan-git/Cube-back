package blog

import (
	"Cube-back/redis"
	"encoding/json"
	"strconv"
)

func blogSendRedis(id int64, blog Blog) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", blog.CubeId)
	b["id"] = id
	b["name"] = userName
	b["cube_id"] = blog.CubeId
	b["cover"] = blog.Cover
	b["title"] = blog.Title
	b["content"] = blog.Content
	b["text"] = blog.Text
	b["image"] = blog.Image
	b["date"] = blog.Date
	b["love"] = blog.Love
	b["comment"] = blog.Comment
	b["collect"] = blog.Collect
	b["view"] = blog.View
	b["label"] = blog.Label
	b["label_type"] = blog.LabelType
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	blogRedisLeftPush(blog.Label, blog.LabelType, redisValue, blog.CubeId)
}

func blogRedisLeftPush(label, labelType, redisString, cubeId string) {
	if label != "" {
		if labelType == "all" {
			redis.LPush("blog_"+label+"_all_new", redisString)
		} else {
			redis.LPush("blog_"+label+"_all_new", redisString)
			redis.LPush("blog_"+labelType+"_new", redisString)
		}
	}
	redis.LPush("blog_new", redisString)
	redis.LPush("profile_blog_"+cubeId, redisString)
}

func blodRedisGet(mode, page, label, labeltype string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var key string
	if labeltype == "" {
		key = mode
	} else if labeltype == "all" {
		key = label + "_all_" + mode
	} else {
		key = labeltype + "_" + mode
	}
	var t = redis.LRange("blog_"+key, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen("blog_" + key)
	return t, l
}

func blogDetailRedisGet(id string) (interface{}, bool) {
	exist := redis.HExists("blog_detail", id)
	if exist {
		var dataBlock []map[string]interface{}
		var m map[string]interface{}
		value := redis.HGet("blog_detail", id)
		json.Unmarshal([]byte(value), &m)
		dataBlock = append(dataBlock, m)
		return dataBlock, true
	}
	return "", false
}

func blogRedisLock(key, status string) {
	redis.Set(key, status)
}

func blogRedisLockStatus(key string) string {
	status := redis.Get(key)
	return status
}

func userCareRedisGet(cubeId string) []string {
	return redis.HKeys("user_cared_" + cubeId)
}
