package gron

import (
	"Cube-back/database"
	"Cube-back/redis"
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

func cubeInformationUpdate() {
	cmd := `select * from user`
	num, _, pass := database.DBValues(cmd)
	if pass {
		l := strconv.FormatInt(num, 10)
		redis.SetNX("user", l, "None", 0)
	}
	cmd = `select * from blog`
	num, _, pass = database.DBValues(cmd)
	if pass {
		l := strconv.FormatInt(num, 10)
		redis.SetNX("blog_count", l, "None", 0)
	}
}

func cubeViewEachMonth() {
	currentMonth := time.Now().Month().String()
	month := Month[currentMonth]
	view := redis.Get("view")
	redis.Set(month, view)
}

func init() {
	c := gron.New()
	c.AddFunc(gron.Every(1*time.Second), func() {
		cubeInformationUpdate()
		cubeViewEachMonth()
	})
	c.Start()
}
