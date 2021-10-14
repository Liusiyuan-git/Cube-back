package blogcomment

import (
	"Cube-back/models/blog"
	"strconv"
	"time"
)

type BlogComment struct {
	Id      int
	CubeId  string `orm:"index"`
	BlogId  int    `orm:"index"`
	Comment string `orm:"type(text)"`
	Love    int
	Date    string `orm:"index;type(datetime)"`
}

var b = new(blog.Blog)

func (bc *BlogComment) BlogCommentSend(blogid, cubeid, comment string) (string, bool) {
	date := time.Now().Format("2006-01-02 15:04:05")
	commentId, err1 := blogCommentSendDb(bc, blogid, cubeid, comment, date)
	if err1 != nil {
		return "评论出错", false
	}
	blogCommentSendDbRedis(blogid, cubeid, comment, date, commentId)
	return "", true
}

func (bc *BlogComment) BlogCommonLike(commentid, blogid, index, love string) (string, bool) {
	bc.Id, _ = strconv.Atoi(commentid)
	bc.Love, _ = strconv.Atoi(love)
	err := blogCommonLikeDb(bc)
	if err != nil {
		return "未知错误", false
	}
	blogCommonLikeRedis(commentid, blogid, index, love)
	return "", true
}

func (bc *BlogComment) BlogCommonGet(blogid, page string) (interface{}, int64, bool) {
	key := "blog_detail_" + blogid + "_comment_get"
	result, length := BlogCommentRedisGet(blogid, page)
	if length != 0 {
		return result, length, true
	}
	if "true" == blogCommentRedisLockStatus(key) {
		return "", 0, false
	}
	blogCommentRedisLock(key, "true")
	result, length, pass := blogCommentDbGet(blogid)
	blogCommentRedisLock(key, "false")
	if pass {
		return result, length, true
	}
	return "", 0, false
}
