package profile

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
)

func profileBlogDbGet(cubeId string) (interface{}, int64, bool) {
	var key = "profile_blog_" + cubeId
	var cmd = `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id and a.cube_id = ?`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		if len(maps) >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return "", 0, false
}

func userProfileBlogDbGet(cubeId string) (interface{}, bool) {
	data := make(map[string]interface{})
	cmd := `select * from user where id=?`
	_, maps, pass := database.DBValues(cmd, cubeId)
	if !pass {
		return "", false
	} else {
		if maps != nil {
			data["name"] = maps[0]["name"]
			data["image"] = maps[0]["image"]
			data["blog"] = maps[0]["blog"]
			data["talk"] = maps[0]["talk"]
			data["collect"] = maps[0]["collect"]
			data["cared"] = maps[0]["cared"]
			data["care"] = maps[0]["care"]
			data["leaving_message"] = maps[0]["leaving_message"]
			data["message"] = maps[0]["message"]
			bjson, _ := json.Marshal(maps[0])
			redisValue := string(bjson)
			redis.HSet("userProfile", cubeId, redisValue)
			return data, true
		} else {
			return "", false
		}
	}
}
