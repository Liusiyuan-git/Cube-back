package talk

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
)

func talkDbGet(mode string) (interface{}, int64, bool) {
	var key = "talk_" + mode
	if "true" == TalkRedisLockStatus(key+"_get") {
		return "数据更新中，请稍后再试", 0, false
	}
	TalkRedisLock(key+"_get", "true")
	var cmd = `select a.id, a.cube_id, a.text, a.date, a.love, a.images, a.comment, b.name FROM talk a inner join user b on a.cube_id = b.cube_id`
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
		if len(maps) >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return "", 0, false
}

func talkDbCmdModeSet(cmd, mode string) string {
	switch mode {
	case "new":
		cmd += " order by id desc"
	case "hot":
		cmd += " order by a.love desc"
	}
	return cmd
}
