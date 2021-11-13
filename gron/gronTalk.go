package gron

import (
	"Cube-back/database"
	"Cube-back/elasticsearch"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
)

func cubeTalkNewUpdate() {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.comment, a.images, b.image as user_image, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by id desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		cubeTalkDataSet(num, maps, "new")
		cubeTalkEsSet(int(num), maps)
	}
}

func cubeTalkDataSet(num int64, maps []orm.Params, mode string) {
	key := "talk_" + mode
	if num != 0 {
		l := redis.LLen(key)
		var i int64
		if l <= num {
			for i = 0; i < num; i++ {
				bjson, _ := json.Marshal(maps[i])
				redisValue := string(bjson)
				if i+1 > l {
					redis.RPush(key, redisValue)
				} else {
					redis.LSet(key, i, redisValue)
				}
				if mode == "new" {
					go cubeTalkCommentUpdate(fmt.Sprintf("%v", maps[i]["id"]))
				}
			}
		} else {
			for i = 0; i < num; i++ {
				bjson, _ := json.Marshal(maps[i])
				redisValue := string(bjson)
				redis.LSet(key, i, redisValue)
				if mode == "new" {
					go cubeTalkCommentUpdate(fmt.Sprintf("%v", maps[i]["id"]))
				}
			}
			redis.LTrim(key, 0, num-1)
		}
	} else {
		redis.LTrim(key, 1, 0)
	}
}

func cubeTalkCommentUpdate(talkid string) {
	cmd := `SELECT a.id, a.comment, a.date, a.cube_id, b.image as user_image, b.name FROM talk_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.talk_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, talkid)
	key := "talk_comment_" + talkid
	if pass {
		if num != 0 {
			var l = redis.LLen(key)
			var i int64
			if l <= num {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(maps[i])
					redisValue := string(bjson)
					if i+1 > l {
						redis.RPush(key, redisValue)
					} else {
						redis.LSet(key, i, redisValue)
					}
				}
			} else {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(maps[i])
					redisValue := string(bjson)
					redis.LSet(key, i, redisValue)
				}
				redis.LTrim(key, 0, num-1)
			}
		} else {
			redis.LTrim(key, 1, 0)
		}
	}
}

func cubeTalkHotUpdate() {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.comment, a.images, b.image as user_image, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by love desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		cubeTalkDataSet(num, maps, "hot")
	}
}

func cubeTalkDetailClean() {
	cmd := `select * from talk`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		l := redis.HGetAll("talk_detail")
		if l != nil {
			for key, _ := range l {
				var keyExits = false
				for _, item := range maps {
					if key == item["id"] {
						keyExits = true
						break
					}
				}
				if !keyExits {
					redis.HDel("talk_detail", key)
					redis.Del("talk_" + key + "_comment_get")
					redis.Del("talk_comment_" + key)
				}
			}
		}
	}
}

func cubeTalkEsSet(num int, maps []orm.Params) {
	EsLen, EsMaps := elasticsearch.Client.SearchAll("blog")
	if num >= EsLen {
		for _, item := range maps {
			var box = map[string]interface{}{}
			box["images"] = item["images"].(string)
			box["name"] = item["name"].(string)
			box["text"] = item["text"].(string)
			box["user_image"] = item["user_image"].(string)
			box["index"], _ = strconv.Atoi(item["id"].(string))
			box["date"] = item["date"].(string)
			box["cube_id"] = item["cube_id"].(string)
			bjson, _ := json.Marshal(box)
			redisValue := string(bjson)
			elasticsearch.Client.Create("talk", redisValue, box["index"].(int))
		}
	} else {
		for index, item := range EsMaps {
			if (index + 1) <= num {
				var box = map[string]interface{}{}
				box["images"] = maps[index]["images"].(string)
				box["user_image"] = maps[index]["user_image"].(string)
				box["name"] = maps[index]["name"].(string)
				box["text"] = maps[index]["text"].(string)
				box["index"], _ = strconv.Atoi(maps[index]["index"].(string))
				box["date"] = maps[index]["date"].(string)
				box["cube_id"] = maps[index]["cube_id"].(string)
				bjson, _ := json.Marshal(box)
				redisValue := string(bjson)
				elasticsearch.Client.Create("talk", redisValue, box["index"].(int))
			} else {
				DocumentId := item.(map[string]interface{})["_id"].(string)
				elasticsearch.Client.Delete("talk", DocumentId)
			}
		}
	}
}
