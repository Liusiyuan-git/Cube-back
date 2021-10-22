package talkcomment

import (
	"Cube-back/database"
	"Cube-back/models/talk"
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

func (b *TalkComment) TalkCommentSend(talkid, cubeid, comment string) (string, bool) {
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
	t := new(talk.Talk)
	redis.HIncrBy("talk_like_and_comment", talkid+"_comment", 1)
	t.Id, _ = strconv.Atoi(talkid)
	t.Comment, _ = strconv.Atoi(redis.HGet("talk_like_and_comment", talkid+"_comment"))
	_, err = database.Update(t, "comment")
	if err != nil {
		return "未知錯誤", false
	}
	TalkCommentRedisSend(talkCommentId, talkid, *b)
	return "", true
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
	talkDetailRedisSet(talkId)
	if pass {
		return result, length, true
	}
	return "", 0, false
}

func (b *TalkComment) TalkCommentDelete(talkcommentid, cubeid, talkid, commentCount, index string) (string, bool) {
	cmd := "DELETE FROM talk_comment where id=? and cube_id=?"
	_, _, pass := database.DBValues(cmd, talkcommentid, cubeid)
	if !pass {
		return "删除失败", false
	}
	talkCommentDeleteRedisUpdate(talkid, index)
	t := new(talk.Talk)
	t.Id, _ = strconv.Atoi(talkid)
	t.Comment, _ = strconv.Atoi(commentCount)
	_, err := database.Update(t, "comment")
	if err != nil {
		return "未知错误", false
	}
	return "", true
}
