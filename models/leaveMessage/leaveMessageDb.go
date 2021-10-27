package leaveMessage

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
	"time"
)

func leaveDbSet(l *LeaveMessage, cubeId, leaveId, text string) error {
	l.Id = 0
	l.CubeId = cubeId
	l.LeaveId = leaveId
	l.Text = text
	l.Date = leaveTime()
	_, err := database.Insert(l)
	if err != nil {
		log.Error(err)
	}
	return err
}

func leaveMessageDbSet(cubeId, leaveId, text string) (*message.Message, error) {
	m := new(message.Message)
	m.CubeId = cubeId
	m.SendId = leaveId
	m.Text = text
	m.Message = 1
	m.Date = time.Now().Format("2006-01-02 15:04:05")
	_, err := database.Insert(m)
	if err != nil {
		log.Error(err)
	}
	return m, err
}

func leaveDbGet(cubeId string) (interface{}, int64, bool) {
	key := "user_leave_" + cubeId
	var cmd = `select a.id, a.leave_id, a.text, a.date, b.image, b.name FROM leave_message a inner join user b on a.cube_id = b.cube_id and a.cube_id=? order by a.id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if pass {
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
	return "", 0, false
}
