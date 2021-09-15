package gron

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/roylee0704/gron"
	"time"
)

var labelType = map[string][]string{
	"all":        []string{},
	"python":     []string{},
	"go":         []string{},
	"java":       []string{},
	"javaScript": []string{},
	"c++":        []string{},
	"c":          []string{},
	"redis":      []string{},
	"rabbitmq":   []string{},
	"docker":     []string{},
	"kubernetes": []string{},
	"mysql":      []string{},
	"live":       []string{},
}

var label = map[string][]string{
	"language":       []string{"python", "go", "java", "javaScript", "c++", "c"},
	"middleware":     []string{"redis", "rabbitmq"},
	"virtualization": []string{"docker", "kubernetes"},
	"database":       []string{"mysql"},
	"other":          []string{"live"},
}

func cubeBlogNewUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.text, a.content, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		cubeBlogDataSplit(maps)
		cubeBlogDataSet(num, labelType["all"], "new")
		cubeBlogDataFilterSet("new")
	}
}

func cubeBlogDataFilterSet(mode string) {
	for key, _ := range label {
		var all []string
		for _, item := range label[key] {
			cubeBlogDataSet(int64(len(labelType[item])), labelType[item], item+"_"+mode)
			all = append(all, labelType[item]...)
		}
		cubeBlogDataSet(int64(len(all)), all, key+"_all_"+mode)
	}
}

func cubeBlogDataSet(num int64, maps []string, mode string) {
	key := "blog_" + mode
	l := redis.LLen(key)
	var i int64
	if l <= num {
		for i = 0; i < num; i++ {
			if i+1 > l {
				redis.RPush(key, maps[i])
			} else {
				redis.LSet(key, i, maps[i])
			}
		}
	} else {
		for i = 0; i < num; i++ {
			redis.LSet(key, i, maps[i])
		}
		redis.LTrim(key, 0, num-1)
	}
}

func cubeBlogDataSplit(maps []orm.Params) {
	labelType["all"] = []string{}
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["label_type"])
		labelType[key] = []string{}
		_, ok := labelType[key]
		var s string
		if ok {
			s = dataConvertToString(item)
			labelType[key] = append(labelType[key], s)
		}
		labelType["all"] = append(labelType["all"], s)
		go cubeBlogDetailUpdate(item)
		go cubeBlogCommentUpdate(fmt.Sprintf("%v", item["id"]))
	}
}

func cubeBlogHotDataSplit(maps []orm.Params) {
	labelType["all"] = []string{}
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["label_type"])
		labelType[key] = []string{}
		_, ok := labelType[key]
		var s string
		if ok {
			s = dataConvertToString(item)
			labelType[key] = append(labelType[key], s)
		}
		labelType["all"] = append(labelType["all"], s)
	}
}

func dataConvertToString(value interface{}) string {
	bjson, _ := json.Marshal(value)
	return string(bjson)
}

func cubeBlogDetailClean() {
	cmd := `select * from blog`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		l := redis.HGetAll("blog_detail")
		if l != nil {
			for key, _ := range l {
				var keyExits = false
				for _, item := range maps {
					if key == item["id"] {
						keyExits = true
						break
					}
				}
				if !keyExits {
					redis.HDel("blog_detail", key)
					redis.Del("blog_detail_" + key + "_get")
					redis.Del("blog_detail_" + key + "_comment_get")
					redis.Del("blog_detail_comment_" + key)
				}
			}
		}
	}
}

func cubeBlogHotUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id order by a.love desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		cubeBlogHotDataSplit(maps)
		cubeBlogDataSet(num, labelType["all"], "hot")
		cubeBlogDataFilterSet("hot")
	}
}

func cubeBlogDetailUpdate(blogDetail orm.Params) {
	id := fmt.Sprintf("%v", blogDetail["id"])
	bjson, _ := json.Marshal(blogDetail)
	redis.HSet("blog_detail", id, string(bjson))
}

func cubeBlogCommentUpdate(blogid string) {
	cmd := `SELECT a.id, a.cube_id, a.comment, a.date, a.love, b.name FROM blog_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.blog_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, blogid)
	key := "blog_detail_comment_" + blogid
	if num != 0 && pass {
		var l = redis.LLen(key)
		var i int64
		if l <= num {
			for i = 0; i < num; i++ {
				bjson, _ := json.Marshal(maps[i])
				redisValue := string(bjson)
				if i+1 > l {
					redis.RPush(key, redisValue)
				} else {
					redis.LSet(key, i, redisValue)
				}
			}
		} else {
			for i = 0; i < num; i++ {
				bjson, _ := json.Marshal(maps[i])
				redisValue := string(bjson)
				redis.LSet(key, i, redisValue)
			}
			redis.LTrim(key, 0, num-1)
		}
	}
}

func cubeTalkNewUpdate() {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.comment, a.images, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by id desc`
	num, maps, pass := database.DBValues(cmd)
	if num != 0 && pass {
		cubeTalkDataSet(num, maps, "new")
	}
}

func cubeTalkHotUpdate() {
	cmd := `select a.id, a.cube_id, a.text, a.date, a.love, a.comment, a.images, b.name FROM talk a inner join user b on a.cube_id = b.cube_id order by love desc`
	num, maps, pass := database.DBValues(cmd)
	if num != 0 && pass {
		cubeTalkDataSet(num, maps, "hot")
	}
}

func cubeTalkDataSet(num int64, maps []orm.Params, mode string) {
	key := "talk_" + mode
	l := redis.LLen(key)
	var i int64
	if l <= num {
		for i = 0; i < num; i++ {
			bjson, _ := json.Marshal(maps[i])
			redisValue := string(bjson)
			if i+1 > l {
				redis.RPush(key, redisValue)
			} else {
				redis.LSet(key, i, redisValue)
			}
			if mode == "new" {
				go cubeTalkCommentUpdate(fmt.Sprintf("%v", maps[i]["id"]))
			}
		}
	} else {
		for i = 0; i < num; i++ {
			bjson, _ := json.Marshal(maps[i])
			redisValue := string(bjson)
			redis.LSet(key, i, redisValue)
			if mode == "new" {
				go cubeTalkCommentUpdate(fmt.Sprintf("%v", maps[i]["id"]))
			}
		}
		redis.LTrim(key, 0, num-1)
	}
}

func cubeTalkCommentUpdate(talkid string) {
	cmd := `SELECT a.id, a.comment, a.date, a.cube_id, b.name FROM talk_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.talk_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, talkid)
	key := "talk_comment_" + talkid
	if num != 0 && pass {
		var l = redis.LLen(key)
		var i int64
		if l <= num {
			for i = 0; i < num; i++ {
				bjson, _ := json.Marshal(maps[i])
				redisValue := string(bjson)
				if i+1 > l {
					redis.RPush(key, redisValue)
				} else {
					redis.LSet(key, i, redisValue)
				}
			}
		} else {
			for i = 0; i < num; i++ {
				bjson, _ := json.Marshal(maps[i])
				redisValue := string(bjson)
				redis.LSet(key, i, redisValue)
			}
			redis.LTrim(key, 0, num-1)
		}
	}
}

func cubeTalkDetailClean() {
	cmd := `select * from talk`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		l := redis.HGetAll("talk_detail")
		if l != nil {
			for key, _ := range l {
				var keyExits = false
				for _, item := range maps {
					if key == item["id"] {
						keyExits = true
						break
					}
				}
				if !keyExits {
					redis.HDel("talk_detail", key)
					redis.Del("talk_" + key + "_comment_get")
					redis.Del("talk_comment_" + key)
				}
			}
		}
	}
}

func init() {
	c := gron.New()
	c.AddFunc(gron.Every(3*time.Second), func() {
		cubeTalkNewUpdate()
		cubeTalkHotUpdate()
	})
	c.AddFunc(gron.Every(3600*time.Second), func() {
		cubeBlogNewUpdate()
		cubeBlogHotUpdate()
		cubeBlogDetailClean()
	})
	c.AddFunc(gron.Every(86400*time.Second), func() {
		cubeBlogDetailClean()
		cubeTalkDetailClean()
	})
	c.Start()
}
