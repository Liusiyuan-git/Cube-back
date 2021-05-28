package redis

import (
	"Cube-back/log"
	"Cube-back/models/common/configure"
	"fmt"
	"github.com/go-redis/redis"
	"reflect"
	"time"
)

type Conf struct {
	RedisIp       string
	RedisPort     string
	RedisPassword string
	Db            int
}

const (
	ip   = "81.68.121.120"
	port = "6379"
	pwd  = "201020120402ssS~"
	db   = 0
)

var client *redis.Client

func timeSet(timeType string, timeLeft time.Duration) time.Duration {
	var t time.Duration
	switch timeType {
	case "Hour":
		t = time.Hour * timeLeft
	case "Minute":
		t = time.Minute * timeLeft
	case "Second":
		t = time.Second * timeLeft
	}
	return t
}

func SetNX(key, value, timeType string, timeLeft time.Duration) {
	t := timeSet(timeType, timeLeft)
	fmt.Print(t)
	status := client.Set(key, value, time.Minute*2).Val()
	println(status)
}

func HSet(key, field, value string) {
	err := client.HSet(key, field, value).Val()
	log.Info(err)
}

func HGet(key, field string) string {
	val := client.HGet(key, field).Val()
	return val
}

func Test() {
	var t time.Duration
	t = time.Minute * 2
	fmt.Println(reflect.TypeOf(t))
}

func init() {
	conf := new(Conf)
	configure.Get(&conf)
	client = redis.NewClient(&redis.Options{
		Addr:     conf.RedisIp + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
		DB:       conf.Db,
		PoolSize: 20,
	})
	p, err := client.Ping().Result()
	if err != nil {
		log.Error(err)
		fmt.Println(err)
	} else {
		log.Info("redis: " + p)
		log.Info("redis init successfully")
	}
}
