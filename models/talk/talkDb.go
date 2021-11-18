package talk

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
	"strconv"
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
			return maps[0:10], num, true
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

func talkDeleteDb(images, cubeId, talkId, date string) (string, bool) {
	cmd := "DELETE FROM talk where id=? and cube_id=?"
	_, _, pass := database.DBValues(cmd, talkId, cubeId)
	if !pass {
		return "删除失败", false
	}
	db := new(DeleteTalk)
	db.TalkId, _ = strconv.Atoi(talkId)
	db.Images = images
	db.CubeId = cubeId
	db.Date = date
	database.Insert(db)
	return "", true
}

func talkMessageSendDb(talkid int64, b *Talk) {
	var userTalkText string
	m := new(message.Message)
	caredBox := userCareRedisGet(b.CubeId)
	for _, item := range caredBox {
		m.CubeId = item
		m.SendId = b.CubeId
		m.Text = userTalkText
		m.Talk = 1
		m.TalkId = int(talkid)
		m.Date = b.Date
		messageId, err := database.Insert(m)
		if err != nil {
			log.Error(err)
			continue
		}
		talkMessageSendRedis(item, messageId, talkid, b)
	}
}
