package gron

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
	"github.com/roylee0704/gron"
	"strconv"
	"time"
)

var Month = map[string]string{
	"January":   "1",
	"February":  "2",
	"March":     "3",
	"April":     "4",
	"May":       "5",
	"June":      "6",
	"July":      "7",
	"August":    "8",
	"September": "9",
	"October":   "10",
	"November":  "11",
	"December":  "12"}

func cubeViewEachMonth() {
	currentMonth := time.Now().Month().String()
	month := Month[currentMonth]
	view := redis.Get("view")
	redis.Set(month, view)
}

func cubeBlogNewUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by id desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		l1 := strconv.FormatInt(num, 10)
		redis.Set("blog_count", l1)
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

func cubeUserUpdate() {
	cmd := `select * from user`
	num, _, pass := database.DBValues(cmd)
	if pass {
		l := strconv.FormatInt(num, 10)
		redis.Set("user", l)
	}
}

func init() {
	c := gron.New()
	c.AddFunc(gron.Every(1*time.Second), func() {
		cubeViewEachMonth()
	})
	c.AddFunc(gron.Every(60*time.Second), func() {
		cubeBlogNewUpdate()
		cubeBlogHotUpdate()
	})
	c.AddFunc(gron.Every(86400*time.Second), func() {
		cubeUserUpdate()
	})
	c.Start()
}
