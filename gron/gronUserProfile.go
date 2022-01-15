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
	userProfileMessageUpdate(userMaps)
}

func userProfileInformationUpdate() []orm.Params {
	cmd := `select * from user`
	num, maps, pass := database.DBValues(cmd)
	if num != 0 && pass {
		userEsUpdate(int(num), maps)
		txpipeline := redis.TxPipeline()
		for _, item := range maps {
			cubeId := fmt.Sprintf("%v", item["cube_id"])
			key := "user_profile_" + cubeId
			txpipeline.HSet(key, "name", fmt.Sprintf("%v", item["name"]))
			txpipeline.HSet(key, "image", fmt.Sprintf("%v", item["image"]))
			txpipeline.HSet(key, "introduce", fmt.Sprintf("%v", item["introduce"]))
		}
		txpipeline.Exec()
		txpipeline.Close()
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

func userTalkClean(userMaps []orm.Params, userTalkBox map[string][]string) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "profile_talk_" + cubeId
		if _, ok := userTalkBox[cubeId]; !ok {
			txpipeline.Del(key)
			txpipeline.HSet("user_profile_"+cubeId, "talk", strconv.Itoa(0))
			userDataBaseTalkUpdate(cubeId, 0)
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func userProfileBlogUpdate(userMaps []orm.Params) {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.date, a.label, a.label_type, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userBlogBox := userSplit(maps)
		userBlogUpdate(userBlogBox)
		userBlogClean(userMaps, userBlogBox)
	}
}

func userBlogClean(userMaps []orm.Params, userBlogBox map[string][]string) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "profile_blog_" + cubeId
		if _, ok := userBlogBox[cubeId]; !ok {
			txpipeline.Del(key)
			txpipeline.HSet("user_profile_"+cubeId, "blog", strconv.Itoa(0))
			userDataBaseBlogUpdate(cubeId, 0)
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func userBlogUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "profile_blog_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		txpipeline := redis.TxPipeline()
		if num != 0 {
			if l <= num {
				for i = 0; i < num; i++ {
					if i+1 > l {
						txpipeline.RPush(key, v[i])
					} else {
						txpipeline.LSet(key, i, v[i])
					}
				}
			} else {
				for i = 0; i < num; i++ {
					txpipeline.LSet(key, i, v[i])
				}
				txpipeline.LTrim(key, 0, num-1)
			}
		} else {
			txpipeline.LTrim(key, 1, 0)
		}
		userDataBaseBlogUpdate(k, len(v))
		txpipeline.HSet("user_profile_"+k, "blog", strconv.Itoa(len(v)))
		txpipeline.Exec()
		txpipeline.Close()
	}
}

func userTalkUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "profile_talk_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		txpipeline := redis.TxPipeline()
		if num != 0 {
			if l <= num {
				for i = 0; i < num; i++ {
					if i+1 > l {
						txpipeline.RPush(key, v[i])
					} else {
						txpipeline.LSet(key, i, v[i])
					}
				}
			} else {
				for i = 0; i < num; i++ {
					txpipeline.LSet(key, i, v[i])
				}
				txpipeline.LTrim(key, 0, num-1)
			}
		} else {
			txpipeline.LTrim(key, 1, 0)
		}
		userDataBaseTalkUpdate(k, len(v))
		txpipeline.HSet("user_profile_"+k, "talk", strconv.Itoa(len(v)))
		txpipeline.Exec()
		txpipeline.Close()
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

func userDataBaseTalkUpdate(cubeId string, length int) {
	b := new(user.User)
	b.CubeId = cubeId
	b.Talk = length
	_, err := database.Update(b, "talk")
	if err != nil {
		log.Error(err)
	}
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
		id := item["collect_cube_id"].(string)
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
	cmd := `SELECT a.cube_id as collect_cube_id, b.id, b.cube_id, b.title, b.cover, b.date, b.title, b.label_type FROM collect a INNER JOIN blog b ON a.blog_id = b.id ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		splitBox := userCollectDateSplit(maps)
		userBlogCollectRedisUpdate(num, splitBox)
		userBlogCollectDbUpdate(splitBox)
		userBlogCollectClean(userMaps, splitBox)
	}
}

func userBlogCollectClean(userMaps []orm.Params, splitBox map[string][]interface{}) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_collect_" + cubeId
		if _, ok := splitBox[cubeId]; !ok {
			txpipeline.Del(key)
			txpipeline.HSet("user_profile_"+cubeId, "collect", "0")
			u := new(user.User)
			u.CubeId = cubeId
			u.Collect = 0
			database.Update(u, "collect")
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
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
		txpipeline := redis.TxPipeline()
		if number != 0 {
			var l = redis.LLen(key)
			var i int64
			if l <= num {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(v[i])
					redisValue := string(bjson)
					if i+1 > l {
						txpipeline.RPush(key, redisValue)
					} else {
						txpipeline.LSet(key, i, redisValue)
					}
				}
			} else {
				for i = 0; i < num; i++ {
					bjson, _ := json.Marshal(v[i])
					redisValue := string(bjson)
					txpipeline.LSet(key, i, redisValue)
				}
				txpipeline.LTrim(key, 0, num-1)
			}
		} else {
			txpipeline.LTrim(key, 1, 0)
		}
		txpipeline.HSet("user_profile_"+k, "collect", fmt.Sprintf("%v", num))
		txpipeline.Exec()
		txpipeline.Close()
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

func userLeaveClean(userMaps []orm.Params, userLeaveBox map[string][]string) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_leave_" + cubeId
		if _, ok := userLeaveBox[cubeId]; !ok {
			txpipeline.Del(key)
			b := new(user.User)
			b.CubeId = cubeId
			b.LeavingMessage = 0
			_, err := database.Update(b, "leaving_message")
			if err != nil {
				log.Error(err)
			}
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func userLeaveUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "user_leave_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		txpipeline := redis.TxPipeline()
		if num != 0 {
			if l <= num {
				for i = 0; i < num; i++ {
					if i+1 > l {
						txpipeline.RPush(key, v[i])
					} else {
						txpipeline.LSet(key, i, v[i])
					}
				}
			} else {
				for i = 0; i < num; i++ {
					txpipeline.LSet(key, i, v[i])
				}
				txpipeline.LTrim(key, 0, num-1)
			}
		} else {
			txpipeline.LTrim(key, 1, 0)
		}
		txpipeline.Exec()
		txpipeline.Close()
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

func userProfileMessageUpdate(userMaps []orm.Params) {
	cmd := `select id, cube_id, send_id, date, text, blog, talk, care, message, blog_comment, talk_comment, blog_id, talk_id FROM message order by id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		userMessageBox := userSplit(maps)
		userMessageUpdate(userMessageBox)
		userMessageClean(userMaps, userMessageBox)
	}
}

func userMessageClean(userMaps []orm.Params, userMessageBox map[string][]string) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_message_" + cubeId
		if _, ok := userMessageBox[cubeId]; !ok {
			txpipeline.Del(key)
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func userMessageUpdate(box map[string][]string) {
	for k, v := range box {
		var key = "user_message_" + k
		var l = redis.LLen(key)
		var num = int64(len(v))
		var i int64
		txpipeline := redis.TxPipeline()
		if num != 0 {
			if l <= num {
				for i = 0; i < num; i++ {
					if i+1 > l {
						txpipeline.RPush(key, v[i])
					} else {
						txpipeline.LSet(key, i, v[i])
					}
				}
			} else {
				for i = 0; i < num; i++ {
					txpipeline.LSet(key, i, v[i])
				}
				txpipeline.LTrim(key, 0, num-1)
			}
		} else {
			txpipeline.LTrim(key, 1, 0)
		}
		txpipeline.Exec()
		txpipeline.Close()
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
		txpipeline := redis.TxPipeline()
		txpipeline.Del(key)
		for _, item := range v {
			txpipeline.HSet(key, item, "1")
		}
		b := new(user.User)
		b.CubeId = k
		b.Care = len(v)
		_, err := database.Update(b, "care")
		if err != nil {
			log.Error(err)
		}
		txpipeline.HSet("user_profile_"+k, "care", fmt.Sprintf("%v", len(v)))
		txpipeline.Exec()
		txpipeline.Close()
	}
}

func userCaredUpdate(box map[string][]string) {
	for k, v := range box {
		key := "user_cared_" + k
		txpipeline := redis.TxPipeline()
		txpipeline.Del(key)
		for _, item := range v {
			txpipeline.HSet(key, item, "1")
		}
		b := new(user.User)
		b.CubeId = k
		b.Cared = len(v)
		_, err := database.Update(b, "cared")
		if err != nil {
			log.Error(err)
		}
		txpipeline.HSet("user_profile_"+k, "cared", fmt.Sprintf("%v", len(v)))
		txpipeline.Exec()
		txpipeline.Close()
	}
}

func userCareClean(userMaps []orm.Params, box map[string][]string) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_care_" + cubeId
		if _, ok := box[cubeId]; !ok {
			txpipeline.Del(key)
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func userCaredClean(userMaps []orm.Params, box map[string][]string) {
	txpipeline := redis.TxPipeline()
	for _, item := range userMaps {
		cubeId := item["cube_id"].(string)
		key := "user_cared_" + cubeId
		if _, ok := box[cubeId]; !ok {
			txpipeline.Del(key)
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
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
				box["index"], _ = strconv.Atoi(maps[index]["id"].(string))
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
