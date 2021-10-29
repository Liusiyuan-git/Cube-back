package gron

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/user"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
)

func userProfileUpdate() {
	userProfileInformationUpdate()
	userProfileBlogUpdate()
	userProfileTalkUpdate()
	userProfileCollectUpdate()
}

func userProfileInformationUpdate() {
	cmd := `select * from user`
	num, maps, pass := database.DBValues(cmd)
	if num != 0 && pass {
		for _, item := range maps {
			cubeId := fmt.Sprintf("%v", item["cube_id"])
			key := "user_profile_" + cubeId
			redis.HSet(key, "name", fmt.Sprintf("%v", item["name"]))
			redis.HSet(key, "image", fmt.Sprintf("%v", item["image"]))
			redis.HSet(key, "introduce", fmt.Sprintf("%v", item["introduce"]))
		}
	}
}

func userProfileTalkUpdate() {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.images, a.comment, b.image as user_image, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userTalkBox := userSplit(maps)
		userTalkUpdate(userTalkBox)
	}
}

func userProfileBlogUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.text, a.content, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userBlogBox := userSplit(maps)
		userBlogUpdate(userBlogBox)
	}
}

func userBlogUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "profile_blog_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		if l <= num {
			for i = 0; i < num; i++ {
				if i+1 > l {
					redis.RPush(key, v[i])
				} else {
					redis.LSet(key, i, v[i])
				}
			}
		} else {
			for i = 0; i < num; i++ {
				redis.LSet(key, i, v[i])
			}
			redis.LTrim(key, 0, num-1)
		}
		userDataBaseBlogUpdate(k, len(v))
		userRedisBlogUpdate(k, len(v))
	}
}

func userTalkUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "profile_talk_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		if l <= num {
			for i = 0; i < num; i++ {
				if i+1 > l {
					redis.RPush(key, v[i])
				} else {
					redis.LSet(key, i, v[i])
				}
			}
		} else {
			for i = 0; i < num; i++ {
				redis.LSet(key, i, v[i])
			}
			redis.LTrim(key, 0, num-1)
		}
		userDataBaseTalkUpdate(k, len(v))
		userRedisTalkUpdate(k, len(v))
	}
}

func userDataBaseBlogUpdate(cubeId string, length int) {
	b := new(user.User)
	b.CubeId = cubeId
	b.Blog = length
	_, err := database.Update(b, "blog")
	if err != nil {
		log.Error(err)
	}
}

func userRedisBlogUpdate(cubeId string, length int) {
	redis.HSet("user_profile_"+cubeId, "blog", strconv.Itoa(length))
}

func userDataBaseTalkUpdate(cubeId string, length int) {
	b := new(user.User)
	b.CubeId = cubeId
	b.Talk = length
	_, err := database.Update(b, "talk")
	if err != nil {
		log.Error(err)
	}
}

func userRedisTalkUpdate(cubeId string, length int) {
	redis.HSet("user_profile_"+cubeId, "talk", strconv.Itoa(length))
}

func userSplit(maps []orm.Params) map[string][]string {
	var box = map[string][]string{}
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["cube_id"])
		_, ok := box[key]
		if !ok {
			box[key] = []string{}
		}
		box[key] = append(box[key], dataConvertToString(item))
	}
	return box
}

func userCollectDateSplit(maps []orm.Params) map[string][]interface{} {
	var splitBox = map[string][]interface{}{}
	for _, item := range maps {
		id := item["cube_id"].(string)
		_, ok := splitBox[id]
		if !ok {
			splitBox[id] = []interface{}{item}
		} else {
			splitBox[id] = append(splitBox[id], item)
		}
	}
	return splitBox
}

func userProfileCollectUpdate() {
	cmd := `SELECT b.id, b.cube_id, b.title, b.cover, b.date, b.title, b.label_type FROM collect a INNER JOIN blog b ON a.blog_id = b.id ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd)
	if !pass {
		splitBox := userCollectDateSplit(maps)
		userBlogCollectRedisUpdate(num, splitBox)
		userBlogCollectDbUpdate(splitBox)
	}
}

func userBlogCollectDbUpdate(splitBox map[string][]interface{}) {
	for k, v := range splitBox {
		u := new(user.User)
		u.CubeId = k
		u.Collect = len(v)
		database.Update(u, "collect")
	}
}

func userBlogCollectRedisUpdate(number int64, splitBox map[string][]interface{}) {
	for k, v := range splitBox {
		key := "user_collect_" + k
		var num = int64(len(v))
		if number != 0 {
			var l = redis.LLen(key)
			var i int64
			if l <= num {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(v)
					redisValue := string(bjson)
					if i+1 > l {
						redis.RPush(key, redisValue)
					} else {
						redis.LSet(key, i, redisValue)
					}
				}
			} else {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(v)
					redisValue := string(bjson)
					redis.LSet(key, i, redisValue)
				}
				redis.LTrim(key, 0, num-1)
			}
		} else {
			redis.LTrim(key, 1, 0)
		}
		redis.HSet("user_profile_"+k, "collect", fmt.Sprintf("%v", num))
	}
}
