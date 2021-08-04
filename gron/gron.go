package gron

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
	"github.com/roylee0704/gron"
	"time"
)

func cubeBlogNewUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		l2 := redis.LLen("blog_new")
		for index, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			if int64(index) <= l2 {
				redis.LSet("blog_new", int64(index), redisValue)
			} else {
				redis.RPush("blog_new", redisValue)
			}
		}
	}
}

func cubeBlogHotUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by a.love desc limit 10`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		for index, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.LSet("blog_hot", int64(index), redisValue)
		}
	}
}

func init() {
	c := gron.New()
	c.AddFunc(gron.Every(60*time.Second), func() {
		cubeBlogNewUpdate()
		cubeBlogHotUpdate()
	})
	c.Start()
}
