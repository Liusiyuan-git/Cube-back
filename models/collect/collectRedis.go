package collect

import (
	"Cube-back/redis"
	"encoding/json"
)

func BlogCollectRedis(cubeid, blogid, cover, date, title, labelType string) {
	redis.HIncrBy("blog_profile_"+blogid, "collect", 1)
	redis.HIncrBy("user_profile_"+cubeid, "collect", 1)
	redis.HSet("blog_profile_"+blogid, cubeid, "1")
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
	redis.LPush("user_collect_"+cubeid, redisValue)
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
