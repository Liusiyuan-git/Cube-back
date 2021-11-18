package talkcomment

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
)

func talkCommentDbGet(talkId string) (interface{}, int64, bool) {
	cmd := `SELECT a.id, a.comment, a.date, a.cube_id, b.image as user_image, b.name FROM talk_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.talk_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, talkId)
	if !pass {
		return "", 0, false
	} else {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush("talk_comment_"+talkId, redisValue)
		}
		if num >= 10 {
			return maps[0:10], num, true
		} else {
			return maps[0:], num, true
		}
	}
}

func talkCommentMessageSendDb(talkCubeId, cubeid string, b *TalkComment) {
	m := new(message.Message)
	m.CubeId = talkCubeId
	m.SendId = cubeid
	m.Text = b.Comment
	m.TalkComment = 1
	m.Date = b.Date
	m.TalkId = b.TalkId
	msgId, err := database.Insert(m)
	if err != nil {
		log.Error(err)
	}
	talkCommentMessageSendRedis(talkCubeId, cubeid, msgId, b)
}
