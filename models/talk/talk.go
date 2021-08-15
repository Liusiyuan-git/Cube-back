package talk

import (
	"Cube-back/database"
	"strconv"
	"time"
)

type Talk struct {
	Id      int
	CubeId  string `orm:"index"`
	Text    string `orm:"type(text)"`
	Date    string `orm:"index;type(datetime)"`
	Love    int
	Comment int
}

func (b *Talk) TalkGet() (interface{}, bool) {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.comment, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by id desc`
	_, maps, pass := database.DBValues(cmd)
	if !pass {
		return "", false
	} else {
		return maps, true
	}
}

func (b *Talk) TalkSend(cubeid, text string) (string, bool) {
	b.Id = 0
	b.CubeId = cubeid
	b.Comment = 0
	b.Text = text
	b.Love = 0
	b.Date = time.Now().Format("2006-01-02 15:04:05")
	_, err := database.Insert(b)
	if err != nil {
		return "发送出错", false
	}
	return "", true
}

func (b *Talk) TalkLike(talkid, like string) (string, bool) {
	b.Id, _ = strconv.Atoi(talkid)
	b.Love, _ = strconv.Atoi(like)
	_, err := database.Update(b, "love")
	if err != nil {
		return "未知錯誤", false
	}
	return "", true
}
