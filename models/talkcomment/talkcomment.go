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

func (b *TalkComment) TalkCommentSend(talkid, cubeid, comment, commentCount string) (string, bool) {
	id, _ := strconv.Atoi(talkid)
	b.Id = 0
	b.TalkId = id
	b.CubeId = cubeid
	b.Comment = comment
	b.Date = time.Now().Format("2006-01-02 15:04:05")
	_, err := database.Insert(b)
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
	return "", true
}

func (b *TalkComment) TalkCommonGet(talkid string) (interface{}, bool) {
	cmd := `SELECT a.id, a.comment, a.date, a.cube_id, b.name FROM talk_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.talk_id = ? ORDER BY a.id DESC`
	_, maps, pass := database.DBValues(cmd, talkid)
	if !pass {
		return "", false
	} else {
		return maps, true
	}
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
