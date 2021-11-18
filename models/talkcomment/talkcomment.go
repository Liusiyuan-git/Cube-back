package talkcomment

import (
	"Cube-back/database"
	"Cube-back/redis"
	"strconv"
	"time"
)

type TalkComment struct {
	Id      int
	CubeId  string `orm:"index"`
	TalkId  int    `orm:"index"`
	Comment string `orm:"type(text)"`
	Date    string `orm:"index;type(datetime)"`
}

func (b *TalkComment) TalkCommentSend(talkid, cubeid, talkCubeId, comment string) (string, bool) {
	id, _ := strconv.Atoi(talkid)
	b.Id = 0
	b.TalkId = id
	b.CubeId = cubeid
	b.Comment = comment
	b.Date = time.Now().Format("2006-01-02 15:04:05")
	talkCommentId, err := database.Insert(b)
	if err != nil {
		return "评论出错", false
	}
	redis.HIncrBy("talk_like_and_comment", talkid+"_comment", 1)
	TalkCommentRedisSend(talkCommentId, talkid, *b)
	if cubeid != talkCubeId {
		go talkCommentMessageSend(talkCubeId, cubeid, b)
	}
	return "", true
}

func talkCommentMessageSend(talkCubeId, cubeid string, b *TalkComment) {
	talkCommentMessageSendDb(talkCubeId, cubeid, b)
}

func (b *TalkComment) TalkCommonGet(talkId, page string) (interface{}, int64, bool) {
	key := "talk_" + talkId + "_comment_get"
	result, length := TalkCommentRedisGet(talkId, page)
	if length != 0 {
		return result, length, true
	}
	if "true" == talkCommentRedisLockStatus(key) {
		return "", 0, false
	}
	talkCommentRedisLock(key, "true")
	result, length, pass := talkCommentDbGet(talkId)
	talkCommentRedisLock(key, "false")
	if pass {
		return result, length, true
	}
	return "", 0, false
}

func (b *TalkComment) TalkCommentDelete(talkcommentid, cubeid, talkid, index string) (string, bool) {
	cmd := "DELETE FROM talk_comment where id=? and cube_id=?"
	_, _, pass := database.DBValues(cmd, talkcommentid, cubeid)
	if !pass {
		return "删除失败", false
	}
	talkCommentDeleteRedisUpdate(talkcommentid, talkid, index)
	return "", true
}
