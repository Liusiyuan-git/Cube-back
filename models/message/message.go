package message

import "encoding/json"

type Message struct {
	Id          int
	CubeId      string
	SendId      string
	Text        string `orm:"type(text)"`
	Date        string `orm:"index;type(datetime)"`
	Blog        int
	Talk        int
	Message     int
	Care        int
	BlogComment int
	TalkComment int
}

func (m *Message) UserMessageGet(cubeId, page string) (interface{}, int64, bool) {
	var dataBlock []map[string]interface{}
	talkData, length := userMessageRedisGet(cubeId, page)
	if len(talkData) == 0 {
		talkDb, length, pass := userMessageDbGet(cubeId)
		return talkDb, length, pass
	}
	for _, item := range talkData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, length, true
}

func (m *Message) MessageProfileGet(cubeId string) interface{} {
	return MessageProfileRedisGet(cubeId)
}
