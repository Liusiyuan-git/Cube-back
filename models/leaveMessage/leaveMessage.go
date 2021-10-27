package leaveMessage

import (
	"Cube-back/rabbitmq"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"time"
)

type LeaveMessage struct {
	Id      int
	CubeId  string `orm:"index"`
	LeaveId string
	Text    string `orm:"type(text)"`
	Date    string `orm:"type(datetime)"`
}

func (l *LeaveMessage) LeaveSet(cubeId, leaveId, text string) bool {
	err1 := leaveDbSet(l, cubeId, leaveId, text)
	m, err2 := leaveMessageDbSet(cubeId, leaveId, text)
	if err1 != nil || err2 != nil {
		return false
	}
	leaveRedisSet(cubeId, leaveId, l)
	leaveMessageRedisSet(m)
	go rabbitmq.MessageQueue.MessageSend(cubeId, fmt.Sprintf("%v", redis.HIncrBy("user_message_profile_"+cubeId, "total", 1)))
	return true
}

func (l *LeaveMessage) LeaveGet(cubeId, page string) (interface{}, int64, bool) {
	var dataBlock []interface{}
	leaveData, length := leaveRedisGet(cubeId, page)
	if length == 0 {
		leaveData, length, pass := leaveDbGet(cubeId)
		return leaveData, length, pass
	}
	for _, item := range leaveData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, length, true
}

func leaveTime() string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t
}
