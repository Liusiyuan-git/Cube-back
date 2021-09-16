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

func userProfileRedisGet(cubeid string) string {
	profile := redis.HGet("userProfile", cubeid)
	return profile
}
