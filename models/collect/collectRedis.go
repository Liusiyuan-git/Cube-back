package collect

import (
	"Cube-back/redis"
	"encoding/json"
	"strconv"
)

func BlogCollectRedis(cubeid, blogid, cover, date, title, labelType string) {
	blogCollectNewUpdate(blogid, cubeid, cover, date, title, labelType)
}

func blogCollectNewUpdate(blogid, cubeid, cover, date, title, labelType string) {
	var maps = map[string]string{}
	maps["id"] = blogid
	maps["cover"] = cover
	maps["date"] = date
	maps["title"] = title
	maps["label_type"] = labelType
	bjson, _ := json.Marshal(maps)
	redisValue := string(bjson)
	txpipeline := redis.TxPipeline()
	txpipeline.LPush("user_collect_"+cubeid, redisValue)
	txpipeline.HIncrBy("blog_profile_"+blogid, "collect", 1)
	txpipeline.HIncrBy("user_profile_"+cubeid, "collect", 1)
	txpipeline.HSet("blog_profile_"+blogid, cubeid, "1")
	txpipeline.Exec()
	txpipeline.Close()
}

func BlogCollectConfirmRedisGet(blogid, cubeid string) bool {
	collect := redis.HGet("blog_profile_"+blogid, cubeid)
	if collect != "" {
		return true
	}
	return false
}

func collectRedisGet(cubeid string) ([]string, int) {
	key := "user_collect_" + cubeid
	var t = redis.LRange(key, 0, 9)
	return t, len(t)
}

func collectDeleteRedis(index, blogId, cubeId string) {
	var key = "user_collect_" + cubeId
	var m map[string]interface{}
	txpipeline := redis.TxPipeline()
	location, _ := strconv.Atoi(index)
	each := redis.LIndex(key, int64(location))
	json.Unmarshal([]byte(each), &m)
	if blogId == m["id"] {
		txpipeline.LRem(key, 1, each)
		txpipeline.HIncrBy("user_profile_"+cubeId, "collect", -1)
		txpipeline.HIncrBy("blog_profile_"+blogId, "collect", -1)
		txpipeline.HDel("blog_profile_"+blogId, cubeId)
	} else {
		blogBox := redis.LRange(key, 0, -1)
		for _, item := range blogBox {
			var s map[string]interface{}
			json.Unmarshal([]byte(item), &s)
			if blogId == s["id"] {
				txpipeline.LRem(key, 1, item)
				txpipeline.HIncrBy("user_profile_"+cubeId, "collect", -1)
				txpipeline.HIncrBy("blog_profile_"+blogId, "collect", -1)
				txpipeline.HDel("blog_profile_"+blogId, cubeId)
				break
			}
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}
