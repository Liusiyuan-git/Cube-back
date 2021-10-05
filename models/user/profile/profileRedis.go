package profile

import (
	"Cube-back/redis"
	"strconv"
)

func profileBlogRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var t = redis.LRange("profile_blog_"+cubeId, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen("profile_blog_" + cubeId)
	return t, l
}

func profileTalkRedisGet(cubeId, page string) ([]string, int64) {
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var t = redis.LRange("profile_talk_"+cubeId, (pageInt-1)*10, (pageInt-1)*10+9)
	var l = redis.LLen("profile_talk_" + cubeId)
	return t, l
}

func userProfileRedisGet(cubeid string) interface{} {
	profile := redis.HMGet("user_profile_"+cubeid, []string{"image", "name", "introduce", "blog", "talk", "collect"})
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
