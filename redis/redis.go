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
	_ = client.Set(key, value, t).Val()
}

func Get(key string) string {
	status, err := client.Get(key).Result()
	if err != nil {
		log.Error(err)
	}
	return status
}

func SetVerificationCode(key, value string) {
	_ = client.Set(key, value, 120000000000).Val()
}

func Set(key, value string) {
	_ = client.Set(key, value, 0).Val()
}

func Lock(key, value string) {
	_ = client.Set(key, value, 60000000000).Val()
}

func Del(key string) {
	_ = client.Del(key).Val()
}

func Exists(key string) int64 {
	status, err := client.Exists(key).Result()
	if err != nil {
		log.Error(err)
	}
	return status
}

func HSet(key, field, value string) {
	client.HSet(key, field, value).Val()
}

func HIncrBy(key, field string, num int64) int64 {
	num, err := client.HIncrBy(key, field, num).Result()
	if err != nil {
		log.Error(err)
	}
	return num
}

func HMGet(key string, fields []string) []interface{} {
	val, err := client.HMGet(key, fields...).Result()
	if err != nil {
		log.Error(err)
	}
	return val
}

func HMSet(key string, fields map[string]interface{}) {
	_, err := client.HMSet(key, fields).Result()
	if err != nil {
		log.Error(err)
	}
}

func HGetAll(key string) map[string]string {
	val, err := client.HGetAll(key).Result()
	if err != nil {
		log.Error(err)
	}
	return val
}

func HKeys(key string) []string {
	val, err := client.HKeys(key).Result()
	if err != nil {
		log.Error(err)
	}
	return val
}

func HDel(key, field string) {
	_, err := client.HDel(key, field).Result()
	if err != nil {
		log.Error(err)
	}
}

func HGet(key, field string) string {
	val, err := client.HGet(key, field).Result()
	if err != nil {
		log.Error(err)
	}
	return val
}

func HExists(key, field string) bool {
	val, err := client.HExists(key, field).Result()
	if err != nil {
		log.Error(err)
	}
	return val
}

func LPush(key string, value string) {
	_, err := client.LPush(key, value).Result()
	if err != nil {
		log.Error(err)
	}
}

func LPop(key string) {
	_, err := client.LPop(key).Result()
	if err != nil {
		log.Error(err)
	}
}

func LLen(key string) int64 {
	length, err := client.LLen(key).Result()
	if err != nil {
		log.Error(err)
	}
	return length
}

func RPush(key string, value string) {
	_, err := client.RPush(key, value).Result()
	if err != nil {
		log.Error(err)
	}
}

func LRange(key string, start, stop int64) []string {
	result, err := client.LRange(key, start, stop).Result()
	if err != nil {
		log.Error(err)
	}
	return result
}

func LIndex(key string, index int64) string {
	result, err := client.LIndex(key, index).Result()
	if err != nil {
		log.Error(err)
	}
	return result
}

func LRem(key, value string) {
	_, err := client.LRem(key, 1, value).Result()
	if err != nil {
		log.Error(err)
	}
}

func LSet(key string, index int64, value string) {
	_, err := client.LSet(key, index, value).Result()
	if err != nil {
		log.Error(err)
	}
}

func LTrim(key string, start, stop int64) {
	_, err := client.LTrim(key, start, stop).Result()
	if err != nil {
		log.Error(err)
	}
}

func Pipeline() redis.Pipeliner {
	pipeline := client.Pipeline()
	return pipeline
}
func TxPipeline() redis.Pipeliner {
	txpipeline := client.TxPipeline()
	return txpipeline
}

func StringStringMapCmd() []*redis.StringStringMapCmd {
	return make([]*redis.StringStringMapCmd, 0)

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
