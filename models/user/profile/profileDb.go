package profile

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/user"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"strconv"
)

func profileBlogDbGet(cubeId string) (interface{}, int64, bool) {
	var key = "profile_blog_" + cubeId
	var cmd = `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id and a.cube_id = ? order by a.id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		if len(maps) >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return "", 0, false
}

func profileTalkDbGet(cubeId string) (interface{}, int64, bool) {
	var key = "profile_talk_" + cubeId
	var cmd = `select a.id, a.cube_id, a.text, a.date, a.love, a.images, a.comment, b.name FROM talk a inner join user b on a.cube_id = b.cube_id and a.cube_id = ? order by a.id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		if len(maps) >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return "", 0, false
}

func UserIntroduceDbSend(cubeId, introduce string) bool {
	u := new(user.User)
	u.Introduce = introduce
	u.CubeId = cubeId
	_, err := database.Update(u, "introduce")
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func UserNameDbSend(cubeId, name string) bool {
	u := new(user.User)
	u.Name = name
	u.CubeId = cubeId
	_, err := database.Update(u, "name")
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func profileCollectDbGet(cubeid string) (interface{}, int64, bool) {
	key := "user_collect_" + cubeid
	cmd := `SELECT b.id, b.cube_id, b.title, b.cover, b.date, b.title, b.label_type FROM collect a INNER JOIN blog b ON a.blog_id = b.id  AND a.cube_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, cubeid)
	if num != 0 && !pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
	}
	redis.HSet("user_profile_"+cubeid, "collect", fmt.Sprintf("%v", num))
	return maps, num, pass
}

func userProfileBlogDbGet(cubeId string) (interface{}, bool) {
	data := make(map[string]interface{})
	cmd := `select * from user where cube_id=?`
	_, maps, pass := database.DBValues(cmd, cubeId)
	if !pass {
		return "", false
	} else {
		if maps != nil {
			data["name"] = maps[0]["name"]
			data["image"] = maps[0]["image"]
			data["blog"] = maps[0]["blog"]
			data["talk"] = maps[0]["talk"]
			data["collect"] = maps[0]["collect"]
			data["cared"] = maps[0]["cared"]
			data["care"] = maps[0]["care"]
			data["leaving_message"] = maps[0]["leaving_message"]
			data["message"] = maps[0]["message"]
			data["introduce"] = maps[0]["introduce"]
			bjson, _ := json.Marshal(maps[0])
			redisValue := string(bjson)
			redis.HSet("userProfile", cubeId, redisValue)
			return data, true
		} else {
			return "", false
		}
	}
}

func userCareDbGet(id string) (interface{}, bool) {
	var careBox = []map[string]string{}
	cmd := `SELECT * from care where care=?`
	num, maps, pass := database.DBValues(cmd, id)
	if !pass {
		if num != 0 {
			for _, item := range maps {
				cubeid := fmt.Sprintf("%v", item["cared"])
				profile := redis.HMGet("user_profile_"+cubeid, []string{"image", "name", "introduce"})
				redis.HSet("user_care_"+id, cubeid, "1")
				careBox = append(careBox, map[string]string{"cube_id": cubeid, "image": fmt.Sprintf("%v", profile[0]), "name": fmt.Sprintf("%v", profile[1]), "introduce": fmt.Sprintf("%v", profile[2])})
			}
		}
		return careBox, true
	}
	return careBox, false
}

func userCareDbSet(id, cubeId string) {
	careUpdate(id)
	caredUpdate(cubeId)
}

func careUpdate(id string) {
	care := redis.HGet("user_profile_"+id, "care")
	u := new(user.User)
	u.CubeId = id
	u.Care, _ = strconv.Atoi(care)
	_, err := database.Update(u, "care")
	if err != nil {
		log.Error(err)
	}
}

func caredUpdate(cubeId string) {
	cared := redis.HGet("user_profile_"+cubeId, "cared")
	u := new(user.User)
	u.CubeId = cubeId
	u.Care, _ = strconv.Atoi(cared)
	_, err := database.Update(u, "care")
	if err != nil {
		log.Error(err)
	}
}

func profileCareDbGet(cubeId string) (interface{}, bool) {
	var careDataBox []map[string]interface{}
	cmd := `SELECT * from care where care=?`
	_, maps, pass := database.DBValues(cmd, cubeId)
	if pass {
		for _, item := range maps {
			caredId := fmt.Sprintf("%v", item["cared"])
			profile := redis.HMGet("user_profile_"+caredId, []string{"name", "image", "introduce"})
			redis.HSet("user_care_"+cubeId, caredId, "1")
			careDataBox = append(careDataBox, map[string]interface{}{"cube_id": caredId, "name": profile[0], "image": profile[1], "introduce": profile[2]})
		}
		return careDataBox, true
	}
	return careDataBox, false
}
