package message

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
)

func userMessageDbGet(cubeId string) ([]orm.Params, int64, bool) {
	var key = "user_message_" + cubeId
	var cmd = `select id, send_id, date, text, blog, talk, care, message, blog_comment, talk_comment, blog_id, talk_id FROM message where cube_id = ? order by id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		txpipeline := redis.TxPipeline()
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			txpipeline.RPush(key, redisValue)
		}
		txpipeline.Exec()
		txpipeline.Close()
		if num >= 10 {
			return maps[0:10], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return []orm.Params{}, num, pass
}

func messageDeleteDb(id string) (string, bool) {
	cmd := "DELETE FROM message where id=?"
	_, _, pass := database.DBValues(cmd, id)
	if !pass {
		return "删除失败", false
	}
	return "", true
}
