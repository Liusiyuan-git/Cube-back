package profile

import (
	"Cube-back/models/message"
	"Cube-back/redis"
	"encoding/json"
	"math"
	"strconv"
)

func profileBlogRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var l = redis.LLen("profile_blog_" + cubeId)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var t = redis.LRange("profile_blog_"+cubeId, (pageInt-1)*10, (pageInt-1)*10+9)
	return t, l
}

func profileTalkRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var l = redis.LLen("profile_talk_" + cubeId)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var t = redis.LRange("profile_talk_"+cubeId, (pageInt-1)*10, (pageInt-1)*10+9)
	return t, l
}

func userProfileRedisGet(cubeid string) interface{} {
	profile := redis.HMGet("user_profile_"+cubeid, []string{"image", "name", "introduce", "blog", "talk", "collect", "cared", "care"})
	return profile
}

func UserIntroduceRedisSend(cubeId, introduce string) {
	redis.HSet("user_profile_"+cubeId, "introduce", introduce)
}

func UserNameRedisSend(cubeId, name string) {
	redis.HSet("user_profile_"+cubeId, "name", name)
}

func SendUserImageRedis(cubeId, image string) {
	redis.HSet("user_profile_"+cubeId, "image", image)
}

func userImageRedisGet(cubeId string) string {
	return redis.HGet("user_profile_"+cubeId, "image")
}

func profileCollectRedisGet(cubeid, page string) ([]string, int64) {
	key := "user_collect_" + cubeid
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var l = redis.LLen(key)
	if pageInt > int64(math.Ceil(float64(l)/10)) {
		pageInt = 1
	}
	var t = redis.LRange(key, (pageInt-1)*10, (pageInt-1)*10+9)
	return t, l
}

func userCareRedisGet(id string) []string {
	return redis.HKeys("user_care_" + id)
}

func userCareRedisSet(id, cubeId string) {
	redis.HSet("user_care_"+id, cubeId, "1")
	redis.HSet("user_cared_"+cubeId, id, "1")
	redis.HIncrBy("user_profile_"+id, "care", 1)
	redis.HIncrBy("user_profile_"+cubeId, "cared", 1)
}

func userCareRedisCancelSet(id, cubeId string) {
	redis.HDel("user_care_"+id, cubeId)
	redis.HDel("user_cared_"+cubeId, id)
	redis.HIncrBy("user_profile_"+id, "care", -1)
	redis.HIncrBy("user_profile_"+cubeId, "cared", -1)
}

func profileCareRedisGet(cubeId string) []map[string]interface{} {
	var careDataBox []map[string]interface{}
	careId := redis.HKeys("user_care_" + cubeId)
	for _, item := range careId {
		if item != "" {
			eachProfile := redis.HMGet("user_profile_"+item, []string{"name", "image", "introduce"})
			careDataBox = append(careDataBox, map[string]interface{}{"cube_id": item, "name": eachProfile[0], "image": eachProfile[1], "introduce": eachProfile[2]})
		}
	}
	return careDataBox
}

func profileCaredRedisGet(cubeId string) []map[string]interface{} {
	var careDataBox []map[string]interface{}
	careId := redis.HKeys("user_cared_" + cubeId)
	for _, item := range careId {
		if item != "" {
			eachProfile := redis.HMGet("user_profile_"+item, []string{"name", "image", "introduce"})
			careDataBox = append(careDataBox, map[string]interface{}{"cube_id": item, "name": eachProfile[0], "image": eachProfile[1], "introduce": eachProfile[2]})
		}
	}
	return careDataBox
}
func userMessageRedisSet(message *message.Message) {
	b := make(map[string]interface{})
	userName := redis.HGet("session", message.SendId)
	userImage := redis.HGet("user_profile_"+message.SendId, "image")
	b["cube_id"] = message.CubeId
	b["send_id"] = message.SendId
	b["date"] = message.Date
	b["text"] = message.Text
	b["care"] = message.Care
	b["name"] = userName
	b["image"] = userImage
	bjson, _ := json.Marshal(b)
	redisValue := string(bjson)
	redis.LPush("user_message_"+message.CubeId, redisValue)
}
