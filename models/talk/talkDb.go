package talk

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
)

func talkDbGet(mode string) (interface{}, int64, bool) {
	var key = "talk_" + mode
	if "true" == TalkRedisLockStatus(key+"_get") {
		return "数据更新中，请稍后再试", 0, false
	}
	TalkRedisLock(key+"_get", "true")
	var cmd = `select a.id, a.cube_id, a.text, a.date, a.love, a.images, a.comment, b.image as user_image, b.name FROM talk a inner join user b on a.cube_id = b.cube_id`
	cmd = talkDbCmdModeSet(cmd, mode)
	num, maps, pass := database.DBValues(cmd)
	TalkRedisLock(key+"_get", "false")
	if num != 0 && pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		redis.Set(key+"_get", "false")
		if num >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	} else if num == 0 {
		return "", 0, true
	}
	return "", 0, false
}

func talkDbCmdModeSet(cmd, mode string) string {
	switch mode {
	case "new":
		cmd += " order by a.id desc"
	case "hot":
		cmd += " order by a.love desc"
	}
	return cmd
}

func talkMessageSendDb(b *Talk) {
	var userTalkText string
	m := new(message.Message)
	caredBox := userCareRedisGet(b.CubeId)
	if len(b.Text) > 30 {
		userTalkText = b.Text[30:]
	} else {
		userTalkText = b.Text
	}
	for _, item := range caredBox {
		m.CubeId = item
		m.SendId = b.CubeId
		m.Text = userTalkText
		m.Talk = 1
		m.Date = b.Date
		_, err := database.Insert(m)
		if err != nil {
			log.Error(err)
		}
	}
}
