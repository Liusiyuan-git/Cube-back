package blog

import (
	"Cube-back/models/collect"
	"Cube-back/rabbitmq"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	Redis "github.com/go-redis/redis"
	"math"
	"strconv"
	"strings"
)

func blogSendRedis(id int64, blog Blog) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", blog.CubeId)
	userImage := redis.HGet("user_profile_"+blog.CubeId, "image")
	b["id"] = strconv.Itoa(int(id))
	b["name"] = userName
	b["cube_id"] = blog.CubeId
	b["cover"] = blog.Cover
	b["title"] = blog.Title
	b["text"] = blog.Text
	b["image"] = blog.Image
	b["date"] = blog.Date
	b["love"] = blog.Love
	b["comment"] = blog.Comment
	b["collect"] = blog.Collect
	b["view"] = blog.View
	b["label"] = blog.Label
	b["label_type"] = blog.LabelType
	b["user_image"] = userImage
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	blogRedisLeftPush(blog.Label, blog.LabelType, redisValue, blog.CubeId)
}

func blogRedisLeftPush(label, labelType, redisString, cubeId string) {
	txpipeline := redis.TxPipeline()
	if label != "" {
		if labelType == "all" {
			txpipeline.LPush("blog_"+label+"_all_new", redisString)
		} else {
			txpipeline.LPush("blog_"+label+"_all_new", redisString)
			txpipeline.LPush("blog_"+labelType+"_new", redisString)
		}
	}
	txpipeline.LPush("blog_new", redisString)
	txpipeline.LPush("profile_blog_"+cubeId, redisString)
	txpipeline.HIncrBy("user_profile_"+cubeId, "blog", 1)
	txpipeline.Exec()
	txpipeline.Close()
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
	var l = redis.LLen("blog_" + key)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var t = redis.LRange("blog_"+key, (pageInt-1)*10, (pageInt-1)*10+9)
	return t, l
}

func blogProfileRedisGet(ids string) interface{} {
	var block []*Redis.SliceCmd
	txpipeline := redis.TxPipeline()
	for _, id := range strings.Split(ids, ";") {
		block = append(block, txpipeline.HMGet("blog_profile_"+id, "love", "comment", "collect", "view"))
	}
	txpipeline.Exec()
	txpipeline.Close()
	var profileBlogBlock [][]interface{}
	for _, item := range block {
		profileBlogBlock = append(profileBlogBlock, item.Val())
	}
	return profileBlogBlock
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

func blogDeleteRedis(label, labelType, index, blogId, cubeId string) {
	c := new(collect.Collect)
	profileBlogRedisClean(index, blogId, cubeId)
	if c.BlogCollectConfirm(blogId, cubeId) {
		c.CollectDelete("0", blogId, cubeId)
	}
	userProfileClean(cubeId)
	go blogRelationsClean(label, labelType, blogId, cubeId)
}

func userProfileClean(cubeId string) {
	redis.HIncrBy("user_profile_"+cubeId, "blog", -1)
}

func blogRelationsClean(label, labelType, blogId, cubeId string) {
	blogDetailRedisClean(blogId)
	blogRedisClean(blogId)
	blogLabelAllRedisNewHotClean(label, blogId, "new")
	blogLabelAllRedisNewHotClean(label, blogId, "hot")
	blogLabelTypeRedisNewHotClean(labelType, blogId, "new")
	blogLabelTypeRedisNewHotClean(labelType, blogId, "hot")
}

func blogLabelAllRedisNewHotClean(label, blogId, mode string) {
	key := "blog_" + label + "_all_" + mode
	blogBox := redis.LRange(key, 0, -1)
	for _, item := range blogBox {
		var s map[string]interface{}
		json.Unmarshal([]byte(item), &s)
		if blogId == s["id"] {
			redis.LRem(key, item)
			break
		}
	}
}

func blogLabelTypeRedisNewHotClean(labelType, blogId, mode string) {
	key := "blog_" + labelType + "_" + mode
	blogBox := redis.LRange(key, 0, -1)
	for _, item := range blogBox {
		var s map[string]interface{}
		json.Unmarshal([]byte(item), &s)
		if blogId == s["id"] {
			redis.LRem(key, item)
			break
		}
	}
}

func blogDetailRedisClean(blogId string) {
	txpipeline := redis.TxPipeline()
	txpipeline.HDel("blog_detail", blogId)
	txpipeline.HDel("blog_message_detail", "cover_"+blogId)
	txpipeline.HDel("blog_message_detail", "title_"+blogId)
	txpipeline.HDel("blog_message_detail", "date_"+blogId)
	txpipeline.HDel("blog_message_detail", "type_"+blogId)
	txpipeline.Del("blog_detail_comment_" + blogId)
	txpipeline.Del("blog_profile_" + blogId)
	txpipeline.Exec()
	txpipeline.Close()
}

func blogRedisClean(blogId string) {
	blogRedisNewHotClean(blogId, "new")
	blogRedisNewHotClean(blogId, "hot")
}

func blogRedisNewHotClean(blogId, mode string) {
	key := "blog_" + mode
	blogBox := redis.LRange(key, 0, -1)
	for _, item := range blogBox {
		var s map[string]interface{}
		json.Unmarshal([]byte(item), &s)
		if blogId == s["id"] {
			redis.LRem(key, item)
			break
		}
	}
}

func profileBlogRedisClean(index, blogId, cubeId string) {
	var key = "profile_blog_" + cubeId
	var m map[string]interface{}
	location, _ := strconv.Atoi(index)
	each := redis.LIndex(key, int64(location))
	json.Unmarshal([]byte(each), &m)
	if blogId == m["id"] {
		redis.LRem(key, each)
	} else {
		blogBox := redis.LRange(key, 0, -1)
		for _, item := range blogBox {
			var s map[string]interface{}
			json.Unmarshal([]byte(item), &s)
			if blogId == s["id"] {
				redis.LRem(key, item)
				break
			}
		}
	}
}

func blogRedisLock(key, status string) {
	redis.Lock(key, status)
}

func blogRedisLockStatus(key string) string {
	status := redis.Get(key)
	return status
}

func userCareRedisGet(cubeId string) []string {
	return redis.HKeys("user_cared_" + cubeId)
}

func blogMessageSendRedis(cubeId string, messageId, blogid int64, t *Blog) {
	b := make(map[string]interface{})
	txpipeline := redis.TxPipeline()
	b["send_id"] = t.CubeId
	b["id"] = strconv.FormatInt(messageId, 10)
	b["date"] = t.Date
	b["blog"] = "1"
	b["blog_id"] = strconv.FormatInt(blogid, 10)
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	execBlock := []interface{}{
		txpipeline.HIncrBy("user_message_profile_"+cubeId, "total", 1),
		txpipeline.LPush("user_message_"+cubeId, redisValue),
		txpipeline.HIncrBy("user_message_profile_"+cubeId, "blog", 1),
		txpipeline.HIncrBy("user_message_profile_"+cubeId, "blog_"+t.CubeId, 1),
	}
	txpipeline.Exec()
	rabbitmq.MessageQueue.MessageSend(cubeId, fmt.Sprintf("%v", execBlock[0].(*Redis.IntCmd).Val()))
	txpipeline.Close()
}

func blogMessageDetailSet(blogid int64, b *Blog) {
	id := strconv.FormatInt(blogid, 10)
	txpipeline := redis.TxPipeline()
	txpipeline.HSet("blog_message_detail", "cover_"+id, b.Cover)
	txpipeline.HSet("blog_message_detail", "title_"+id, b.Title)
	txpipeline.HSet("blog_message_detail", "date_"+id, b.Date)
	txpipeline.HSet("blog_message_detail", "type_"+id, b.LabelType)
	txpipeline.Exec()
	txpipeline.Close()
}
