package message

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
)

func userMessageDbGet(cubeId string) (interface{}, int64, bool) {
	var key = "user_message_" + cubeId
	var cmd = `select a.id, a.send_id, a.date, a.text, a.blog, a.talk, a.care, a.message, a.blog_comment, a.talk_comment, b.image, b.name FROM message a inner join user b on a.send_id = b.cube_id and a.cube_id = ? order by id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		if num >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return "", num, pass
}
