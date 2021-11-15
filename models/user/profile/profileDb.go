package profile

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/care"
	"Cube-back/models/message"
	"Cube-back/models/user"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"time"
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
		if num >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return maps, num, pass
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
		if num >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return maps, num, pass
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
		if num >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return maps, num, pass
}

func userProfileDbGet(cubeId string) (interface{}, bool) {
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
	var cmd = `select a.cared, b.name, b.image, b.introduce FROM care a inner join user b on a.cared = b.cube_id and a.care = ? order by a.id desc`
	num, maps, pass := database.DBValues(cmd, id)
	if num != 0 && !pass {
		for _, item := range maps {
			cubeid := fmt.Sprintf("%v", item["cared"])
			image := fmt.Sprintf("%v", item["image"])
			name := fmt.Sprintf("%v", item["name"])
			introduce := fmt.Sprintf("%v", item["introduce"])
			redis.HSet("user_care_"+id, cubeid, "1")
			careBox = append(careBox, map[string]string{"cube_id": cubeid, "image": image, "name": name, "introduce": introduce})
		}
	}
	return careBox, pass
}

func profileCareDbGet(cubeId string) (interface{}, bool) {
	var careDataBox []map[string]interface{}
	var cmd = `select a.cared, b.name, b.image, b.introduce FROM care a inner join user b on a.cared = b.cube_id and a.care = ? order by a.id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		for _, item := range maps {
			caredId := fmt.Sprintf("%v", item["cared"])
			image := fmt.Sprintf("%v", item["image"])
			name := fmt.Sprintf("%v", item["name"])
			introduce := fmt.Sprintf("%v", item["introduce"])
			redis.HSet("user_care_"+cubeId, caredId, "1")
			careDataBox = append(careDataBox, map[string]interface{}{"cube_id": caredId, "name": name, "image": image, "introduce": introduce})
		}
	}
	return careDataBox, pass
}

func profileCaredDbGet(cubeId string) (interface{}, bool) {
	var careDataBox []map[string]interface{}
	var cmd = `select a.care, b.name, b.image, b.introduce FROM care a inner join user b on a.care = b.cube_id and a.cared = ? order by a.id desc`
	num, maps, pass := database.DBValues(cmd, cubeId)
	if num != 0 && pass {
		for _, item := range maps {
			careId := fmt.Sprintf("%v", item["care"])
			image := fmt.Sprintf("%v", item["image"])
			name := fmt.Sprintf("%v", item["name"])
			introduce := fmt.Sprintf("%v", item["introduce"])
			redis.HSet("user_cared_"+cubeId, careId, "1")
			careDataBox = append(careDataBox, map[string]interface{}{"cube_id": careId, "name": name, "image": image, "introduce": introduce})
		}
	}
	return careDataBox, pass
}

func userCareDbSet(id, cubeId string) error {
	c := new(care.Care)
	c.Care = id
	c.Cared = cubeId
	_, err := database.Insert(c)
	if err != nil {
		log.Error(err)
	}
	return err
}

func userCareMessageDbSet(id, cubeId string) (*message.Message, error) {
	m := new(message.Message)
	m.CubeId = id
	m.SendId = cubeId
	m.Text = "感谢关注！！！ 😄"
	m.Care = 1
	m.Date = time.Now().Format("2006-01-02 15:04:05")
	_, err := database.Insert(m)
	if err != nil {
		log.Error(err)
	}
	return m, err
}
