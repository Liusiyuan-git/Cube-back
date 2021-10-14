package blogcomment

import (
	"Cube-back/database"
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
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush("blog_detail_comment_"+blogid, redisValue)
		}
		redis.HSet("blog_profile_"+blogid, "comment", fmt.Sprintf("%v", num))
		if len(maps) >= 10 {
			return maps[0:9], num, true
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
