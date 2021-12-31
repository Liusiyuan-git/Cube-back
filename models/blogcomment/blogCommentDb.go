package blogcomment

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"strconv"
)

func blogCommentDbGet(blogid string) (interface{}, int64, bool) {
	cmd := `SELECT a.id, a.cube_id, a.comment, a.date, a.love, b.image, b.name FROM blog_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.blog_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, blogid)
	if !pass {
		return "", 0, false
	} else {
		txpipeline := redis.TxPipeline()
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			txpipeline.RPush("blog_detail_comment_"+blogid, redisValue)
		}
		txpipeline.HSet("blog_profile_"+blogid, "comment", fmt.Sprintf("%v", num))
		txpipeline.Exec()
		txpipeline.Close()
		if len(maps) >= 10 {
			return maps[0:10], num, true
		} else {
			return maps[0:], num, true
		}
	}
}

func blogCommonLikeDb(bc *BlogComment) error {
	_, err := database.Update(bc, "love")
	return err
}

func blogCommentSendDb(bc *BlogComment, blogid, cubeid, comment, date string) (int64, error) {
	id, _ := strconv.Atoi(blogid)
	bc.Id = 0
	bc.BlogId = id
	bc.CubeId = cubeid
	bc.Comment = comment
	bc.Date = date
	commentId, err1 := database.Insert(bc)
	return commentId, err1
}

func blogCommentMessageSendDb(cubeid, blogCubeId string, bc *BlogComment) {
	m := new(message.Message)
	m.CubeId = blogCubeId
	m.SendId = cubeid
	m.BlogComment = 1
	m.Text = bc.Comment
	m.BlogId = bc.BlogId
	m.Date = bc.Date
	msgId, err := database.Insert(m)
	if err != nil {
		log.Error(err)
	}
	blogCommentMessageSendRedis(blogCubeId, cubeid, msgId, bc)
}
