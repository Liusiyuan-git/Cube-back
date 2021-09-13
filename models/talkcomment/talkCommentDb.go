package talkcomment

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
)

func talkCommentDbGet(talkId string) (interface{}, int64, bool) {
	cmd := `SELECT a.id, a.comment, a.date, a.cube_id, b.name FROM talk_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.talk_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, talkId)
	if !pass {
		return "", 0, false
	} else {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush("talk_comment_"+talkId, redisValue)
		}
		return maps, num, true
	}
}
