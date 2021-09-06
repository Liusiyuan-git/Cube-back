package talkcomment

import (
	"Cube-back/database"
	"Cube-back/models/talk"
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

func (b *TalkComment) TalkCommentSend(talkid, cubeid, index, comment, commentCount, mode string) (string, bool) {
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
	t.Id, _ = strconv.Atoi(talkid)
	t.Comment, _ = strconv.Atoi(commentCount)
	_, err = database.Update(t, "comment")
	if err != nil {
		return "未知錯誤", false
	}
	TalkCommentRedisSend(talkCommentId, talkid, index, mode, commentCount, *b)
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

func (b *TalkComment) TalkCommentDelete(talkcommentid, cubeid, talkid, commentCount string) (string, bool) {
	cmd := "DELETE FROM talk_comment where id=? and cube_id=?"
	_, _, pass := database.DBValues(cmd, talkcommentid, cubeid)
	if !pass {
		return "删除失败", false
	}
	t := new(talk.Talk)
	t.Id, _ = strconv.Atoi(talkid)
	t.Comment, _ = strconv.Atoi(commentCount)
	_, err := database.Update(t, "comment")
	if err != nil {
		return "未知错误", false
	}
	return "", true
}
