package gron

import (
	"Cube-back/database"
	"Cube-back/elasticsearch"
	"Cube-back/log"
	"Cube-back/models/user"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
)

func userProfileUpdate() {
	userMaps := userProfileInformationUpdate()
	userProfileBlogUpdate(userMaps)
	userProfileTalkUpdate(userMaps)
	userProfileCollectUpdate(userMaps)
	userProfileLeaveUpdate(userMaps)
	userProfileCareUpdate(userMaps)
}

func userProfileInformationUpdate() []orm.Params {
	cmd := `select * from user`
	num, maps, pass := database.DBValues(cmd)
	if num != 0 && pass {
		userEsUpdate(int(num), maps)
		for _, item := range maps {
			cubeId := fmt.Sprintf("%v", item["cube_id"])
			key := "user_profile_" + cubeId
			redis.HSet(key, "name", fmt.Sprintf("%v", item["name"]))
			redis.HSet(key, "image", fmt.Sprintf("%v", item["image"]))
			redis.HSet(key, "introduce", fmt.Sprintf("%v", item["introduce"]))
		}
	}
	return maps
}

func userProfileTalkUpdate(userMaps []orm.Params) {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.images, a.comment, b.image as user_image, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userTalkBox := userSplit(maps)
		userTalkUpdate(userTalkBox)
		userTalkClean(userMaps, userTalkBox)
	}
}

func userTalkClean(userMaps []orm.Params, userBlogBox map[string][]string) {
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "profile_talk_" + cubeId
		if _, ok := userBlogBox[cubeId]; !ok {
			redis.Del(key)
			userDataBaseTalkUpdate(cubeId, 0)
			userRedisTalkUpdate(cubeId, 0)
		}
	}
}

func userProfileBlogUpdate(userMaps []orm.Params) {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.text, a.content, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userBlogBox := userSplit(maps)
		userBlogUpdate(userBlogBox)
		userBlogClean(userMaps, userBlogBox)
	}
}

func userBlogClean(userMaps []orm.Params, userBlogBox map[string][]string) {
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "profile_blog_" + cubeId
		if _, ok := userBlogBox[cubeId]; !ok {
			redis.Del(key)
			userDataBaseBlogUpdate(cubeId, 0)
			userRedisBlogUpdate(cubeId, 0)
		}
	}
}

func userBlogUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "profile_blog_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		if num != 0 {
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
		} else {
			redis.LTrim(key, 1, 0)
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
		if num != 0 {
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
		} else {
			redis.LTrim(key, 1, 0)
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

func userProfileCollectUpdate(userMaps []orm.Params) {
	cmd := `SELECT b.id, b.cube_id, b.title, b.cover, b.date, b.title, b.label_type FROM collect a INNER JOIN blog b ON a.blog_id = b.id ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		splitBox := userCollectDateSplit(maps)
		userBlogCollectRedisUpdate(num, splitBox)
		userBlogCollectDbUpdate(splitBox)
		userBlogCollectClean(userMaps, splitBox)
	}
}

func userBlogCollectClean(userMaps []orm.Params, userBlogBox map[string][]interface{}) {
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_collect_" + cubeId
		if _, ok := userBlogBox[cubeId]; !ok {
			redis.Del(key)
			redis.HSet("user_profile_"+cubeId, "collect", "0")
			u := new(user.User)
			u.CubeId = cubeId
			u.Collect = 0
			database.Update(u, "collect")
		}
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
					bjson, _ := json.Marshal(v[i])
					redisValue := string(bjson)
					if i+1 > l {
						redis.RPush(key, redisValue)
					} else {
						redis.LSet(key, i, redisValue)
					}
				}
			} else {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(v[i])
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

func userProfileLeaveUpdate(userMaps []orm.Params) {
	cmd := `select a.id, a.cube_id, a.leave_id, a.text, a.date, b.image, b.name FROM leave_message a inner join user b on a.leave_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userLeaveBox := userSplit(maps)
		userLeaveUpdate(userLeaveBox)
		userLeaveClean(userMaps, userLeaveBox)
	}
}

func userLeaveClean(userMaps []orm.Params, userBlogBox map[string][]string) {
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_leave_" + cubeId
		if _, ok := userBlogBox[cubeId]; !ok {
			redis.Del(key)
			b := new(user.User)
			b.CubeId = cubeId
			b.LeavingMessage = 0
			_, err := database.Update(b, "leaving_message")
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func userLeaveUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "user_leave_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		if num != 0 {
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
		} else {
			redis.LTrim(key, 1, 0)
		}
		userDataBaseLeaveUpdate(k, len(v))
	}
}

func userDataBaseLeaveUpdate(cubeId string, length int) {
	b := new(user.User)
	b.CubeId = cubeId
	b.LeavingMessage = length
	_, err := database.Update(b, "leaving_message")
	if err != nil {
		log.Error(err)
	}
}

func userProfileCareUpdate(userMaps []orm.Params) {
	cmd := `select * FROM care`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userCareBox := userCareSplit(maps)
		userCaredBox := userCaredSplit(maps)
		userCareUpdate(userCareBox)
		userCaredUpdate(userCaredBox)
		userCareClean(userMaps, userCareBox)
		userCaredClean(userMaps, userCaredBox)
	}
}

func userCareSplit(maps []orm.Params) map[string][]string {
	var box = map[string][]string{}
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["care"])
		_, ok := box[key]
		if !ok {
			box[key] = []string{}
		}
		box[key] = append(box[key], item["cared"].(string))
	}
	return box
}

func userCaredSplit(maps []orm.Params) map[string][]string {
	var box = map[string][]string{}
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["cared"])
		_, ok := box[key]
		if !ok {
			box[key] = []string{}
		}
		box[key] = append(box[key], item["care"].(string))
	}
	return box
}

func userCareUpdate(box map[string][]string) {
	for k, v := range box {
		key := "user_care_" + k
		redis.Del(key)
		for _, item := range v {
			redis.HSet(key, item, "1")
		}
		b := new(user.User)
		b.CubeId = k
		b.Care = len(v)
		_, err := database.Update(b, "care")
		if err != nil {
			log.Error(err)
		}
		redis.HSet("user_profile_"+k, "care", fmt.Sprintf("%v", len(v)))
	}
}

func userCaredUpdate(box map[string][]string) {
	for k, v := range box {
		key := "user_cared_" + k
		redis.Del(key)
		for _, item := range v {
			redis.HSet(key, item, "1")
		}
		b := new(user.User)
		b.CubeId = k
		b.Cared = len(v)
		_, err := database.Update(b, "cared")
		if err != nil {
			log.Error(err)
		}
		redis.HSet("user_profile_"+k, "cared", fmt.Sprintf("%v", len(v)))
	}
}

func userCareClean(userMaps []orm.Params, box map[string][]string) {
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_care_" + cubeId
		if _, ok := box[cubeId]; !ok {
			redis.Del(key)
		}
	}
}

func userCaredClean(userMaps []orm.Params, box map[string][]string) {
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_cared_" + cubeId
		if _, ok := box[cubeId]; !ok {
			redis.Del(key)
		}
	}
}

func userEsUpdate(num int, maps []orm.Params) {
	EsLen, EsMaps := elasticsearch.Client.SearchAll("user")
	if num >= EsLen {
		for index, item := range maps {
			var box = map[string]interface{}{}
			box["introduce"] = item["introduce"].(string)
			box["name"] = item["name"].(string)
			box["image"] = item["image"].(string)
			box["index"], _ = strconv.Atoi(item["id"].(string))
			box["cube_id"] = item["cube_id"].(string)
			bjson, _ := json.Marshal(box)
			redisValue := string(bjson)
			elasticsearch.Client.Create("user", redisValue, index)
		}
	} else {
		for index, item := range EsMaps {
			if (index + 1) <= num {
				var box = map[string]interface{}{}
				box["introduce"] = maps[index]["introduce"].(string)
				box["image"] = maps[index]["image"].(string)
				box["name"] = maps[index]["name"].(string)
				box["index"], _ = strconv.Atoi(maps[index]["index"].(string))
				box["cube_id"] = maps[index]["cube_id"].(string)
				bjson, _ := json.Marshal(box)
				redisValue := string(bjson)
				elasticsearch.Client.Create("user", redisValue, index)
			} else {
				DocumentId := item.(map[string]interface{})["_id"].(string)
				elasticsearch.Client.Delete("user", DocumentId)
			}
		}
	}
}
