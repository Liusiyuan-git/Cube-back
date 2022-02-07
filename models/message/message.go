package message

import (
	"Cube-back/models/user"
	"strings"
)

type Message struct {
	Id          int
	BlogId      int
	TalkId      int
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
	_, pass := user.NumberCorrect(cubeId)
	if !pass {
		return []string{}, 0, false
	}
	messageData, length := userMessageRedisGet(cubeId, page)
	if len(messageData) == 0 {
		messageDb, length, pass := userMessageDbGet(cubeId)
		if !pass {
			return []string{}, 0, false
		}
		return userMessageDbRedisDetailGet(messageDb, cubeId), length, pass
	}
	return userMessageRedisDetailGet(messageData), length, true
}

func (m *Message) MessageProfileGet(cubeId string) interface{} {
	return MessageProfileRedisGet(cubeId)
}

func (m *Message) UserMessageClean(cubeId string) {
	UserMessageCleanRedis(cubeId)
}

func (m *Message) MessageProfileUserTalkGet(cubeId, idBox string) (interface{}, bool) {
	var talkIdBox []string
	ids := strings.Split(idBox, ";")
	for _, item := range ids {
		talkIdBox = append(talkIdBox, "talk_"+item)
	}
	return messageProfileUserTalkRedisGet(cubeId, talkIdBox), true
}

func (m *Message) MessageProfileUserTalkClean(id, deleteId string) {
	messageProfileUserTalkRedisClean(id, deleteId)
}

func (m *Message) MessageProfileUserBlogClean(id, deleteId string) {
	messageProfileUserBlogRedisClean(id, deleteId)
}

func (m *Message) MessageProfileUserBlogGet(cubeId, idBox string) (interface{}, bool) {
	var blogIdBox []string
	ids := strings.Split(idBox, ";")
	for _, item := range ids {
		blogIdBox = append(blogIdBox, "blog_"+item)
	}
	return messageProfileUserBlogRedisGet(cubeId, blogIdBox), true
}

func (m *Message) MessageDelete(id, cubeId, index string) (string, bool) {
	msg, pass := user.NumberCorrect(cubeId)
	if !pass {
		return msg, false
	}
	result, pass := messageDeleteDb(id)
	if !pass {
		return result, pass
	}
	messageDeleteRedis(id, cubeId, index)
	return "", true
}
