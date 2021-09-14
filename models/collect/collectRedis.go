package collect

import (
	"Cube-back/redis"
	"encoding/json"
)

func BlogCollectRedis(cudeid, blogid, collect string) {
	blogCollectIdSet(cudeid, blogid)
	result := blogCollectNewUpdate(blogid, collect)
	blogCollectCubeIdSet(cudeid, blogid, result)
}

func blogCollectNewUpdate(blogid, collect string) string {
	var m map[string]interface{}
	value := redis.HGet("blog_detail", blogid)
	json.Unmarshal([]byte(value), &m)
	m["collect"] = collect
	bjson, _ := json.Marshal(m)
	s := string(bjson)
	redis.HSet("blog_detail", blogid, s)
	return s
}

func blogCollectIdSet(cudeid, blogid string) {
	var m = map[string]string{cudeid: "1"}
	bjson, _ := json.Marshal(m)
	redis.HSet("collect_blog_id", blogid, string(bjson))
}

func BlogCollectConfirmRedisGet(id, cubeid string) bool {
	collect := redis.HGet("collect_blog_id", id)
	if collect != "nil" {
		var m map[string]string
		json.Unmarshal([]byte(collect), &m)
		_, ok := m[cubeid]
		return ok
	}
	return false
}

func blogCollectCubeIdSet(cudeid, blogid, result string) {
	collectList := redis.HGet("collect_cube_id", cudeid)
	var m []string
	if collectList != "nil" {
		json.Unmarshal([]byte(collectList), &m)
	}
	bjson, _ := json.Marshal(append(m, result))
	redis.HSet("collect_cube_id", cudeid, string(bjson))
}
