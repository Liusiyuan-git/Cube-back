package gron

import (
	"Cube-back/database"
	"Cube-back/elasticsearch"
	"Cube-back/log"
	"Cube-back/models/blog"
	"Cube-back/redis"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
)

var label = map[string][]string{
	"language":       []string{"python", "go", "java", "javaScript", "c++", "c"},
	"middleware":     []string{"redis", "rabbitmq"},
	"virtualization": []string{"docker", "kubernetes"},
	"database":       []string{"mysql"},
	"other":          []string{"live"},
}

func cubeBlogNewUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.text, a.content, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name, b.image as user_image FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		labelType := labelTypeBuild()
		cubeBlogDataSplit(maps, labelType)
		cubeBlogDataSet(num, labelType["all"], "new")
		cubeBlogDataFilterSet("new", labelType)
		cubeBlogProfileRedisSet(maps)
		//cubeBlogProfileDbSet(maps)
		cubeBlogEsSet(int(num), maps)
	}
}

func labelTypeBuild() map[string][]string {
	return map[string][]string{
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
}

func cubeBlogDataSplit(maps []orm.Params, labelType map[string][]string) {
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["label_type"])
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

func cubeBlogDetailUpdate(blogDetail orm.Params) {
	id := fmt.Sprintf("%v", blogDetail["id"])
	bjson, _ := json.Marshal(blogDetail)
	redis.HSet("blog_detail", id, string(bjson))
}

func cubeBlogCommentUpdate(blogid string) {
	cmd := `SELECT a.id, a.cube_id, a.comment, a.date, a.love, b.image, b.name FROM blog_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.blog_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, blogid)
	key := "blog_detail_comment_" + blogid
	if pass {
		if num != 0 {
			var l = redis.LLen(key)
			if l <= num {
				for i := int64(0); i < num; i++ {
					bjson, _ := json.Marshal(maps[i])
					redisValue := string(bjson)
					if i+1 > l {
						redis.RPush(key, redisValue)
					} else {
						redis.LSet(key, i, redisValue)
					}
				}
			} else {
				for i := int64(0); i < num; i++ {
					bjson, _ := json.Marshal(maps[i])
					redisValue := string(bjson)
					redis.LSet(key, i, redisValue)
				}
				redis.LTrim(key, 0, num-1)
			}
		} else {
			redis.LTrim(key, 1, 0)
		}
		cubeBlogProfileCommentRedisSet(blogid, num)
		cubeBlogProfileCommentDbSet(blogid, num)
	}
}

func cubeBlogProfileCommentRedisSet(blogid string, num int64) {
	redis.HSet("blog_profile_"+blogid, "comment", fmt.Sprintf("%v", num))
}

func cubeBlogProfileCommentDbSet(blogid string, num int64) {
	b := new(blog.Blog)
	b.Id, _ = strconv.Atoi(blogid)
	b.Comment, _ = strconv.Atoi(strconv.FormatInt(num, 10))
	_, err := database.Update(b, "comment")
	if err != nil {
		log.Error(err)
	}
}

func cubeBlogHotUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name, b.image as user_image FROM blog a inner join user b on a.cube_id = b.cube_id order by a.love desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		labelType := labelTypeBuild()
		cubeBlogHotDataSplit(maps, labelType)
		cubeBlogDataSet(num, labelType["all"], "hot")
		cubeBlogDataFilterSet("hot", labelType)
	}
}

func cubeBlogHotDataSplit(maps []orm.Params, labelType map[string][]string) {
	for _, item := range maps {
		key := fmt.Sprintf("%v", item["label_type"])
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

func cubeBlogDataSet(num int64, maps []string, mode string) {
	key := "blog_" + mode
	l := redis.LLen(key)
	if l <= num {
		for i := int64(0); i < num; i++ {
			if i+1 > l {
				redis.RPush(key, maps[i])
			} else {
				redis.LSet(key, i, maps[i])
			}
		}
	} else {
		for i := int64(0); i < num; i++ {
			redis.LSet(key, i, maps[i])
		}
		redis.LTrim(key, 0, num-1)
	}
}

func cubeBlogDataFilterSet(mode string, labelType map[string][]string) {
	for key := range label {
		var all []string
		for _, item := range label[key] {
			length := int64(len(labelType[item]))
			cubeBlogDataSet(length, labelType[item], item+"_"+mode)
			if length != 0 {
				all = append(all, labelType[item]...)
			}
		}
		cubeBlogDataSet(int64(len(all)), all, key+"_all_"+mode)
	}
}

func cubeBlogDetailClean() {
	cmd := `select * from blog`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		l := redis.HGetAll("blog_detail")
		if l != nil {
			for key := range l {
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

func cubeBlogProfileRedisSet(maps []orm.Params) {
	for _, item := range maps {
		id := fmt.Sprintf("%v", item["id"])
		redis.HMSet("blog_profile_"+id, map[string]interface{}{"love": item["love"], "view": item["view"]})
	}
}

//func cubeBlogProfileDbSet(maps []orm.Params) {
//	b := new(blog.Blog)
//	for _, item := range maps {
//		b.Id, _ = strconv.Atoi(item["id"].(string))
//		b.Love, _ = strconv.Atoi(item["love"].(string))
//		b.View, _ = strconv.Atoi(item["view"].(string))
//		_, err := database.Update(b, "love", "comment", "collect", "view")
//		if err != nil {
//			log.Error(err)
//		}
//	}
//}

func cubeBlogEsSet(num int, maps []orm.Params) {
	EsLen, EsMaps := elasticsearch.Client.SearchAll("blog")
	if num >= EsLen {
		for _, item := range maps {
			var box = map[string]interface{}{}
			box["label_type"] = item["label_type"].(string)
			box["name"] = item["name"].(string)
			box["text"] = item["text"].(string)
			box["title"] = item["title"].(string)
			box["user_image"] = item["user_image"].(string)
			box["index"], _ = strconv.Atoi(item["id"].(string))
			box["date"] = item["date"].(string)
			box["cube_id"] = item["cube_id"].(string)
			box["cover"] = item["cover"].(string)
			bjson, _ := json.Marshal(box)
			redisValue := string(bjson)
			elasticsearch.Client.Create("blog", redisValue, box["index"].(int))
		}
	} else {
		for index, item := range EsMaps {
			if (index + 1) <= num {
				var box = map[string]interface{}{}
				box["label_type"] = maps[index]["label_type"].(string)
				box["user_image"] = maps[index]["user_image"].(string)
				box["name"] = maps[index]["name"].(string)
				box["text"] = maps[index]["text"].(string)
				box["title"] = maps[index]["title"].(string)
				box["index"], _ = strconv.Atoi(maps[index]["id"].(string))
				box["date"] = maps[index]["date"].(string)
				box["cube_id"] = maps[index]["cube_id"].(string)
				box["cover"] = maps[index]["cover"].(string)
				bjson, _ := json.Marshal(box)
				redisValue := string(bjson)
				elasticsearch.Client.Create("blog", redisValue, box["index"].(int))
			} else {
				DocumentId := item.(map[string]interface{})["_id"].(string)
				elasticsearch.Client.Delete("blog", DocumentId)
			}
		}
	}
}

func cubeBlogCollectUpdate() {
	cmd := `SELECT * from collect`
	_, maps, pass := database.DBValues(cmd)
	if !pass {
		splitBox := cubeCollectDateSplit(maps)
		cubeBlogCollectRedisUpdate(splitBox)
		cubeBlogCollectDbUpdate(splitBox)
	}
}

func cubeCollectDateSplit(maps []orm.Params) map[string]int {
	var splitBox = map[string]int{}
	for _, item := range maps {
		id := item["blog_id"].(string)
		_, ok := splitBox[id]
		if !ok {
			splitBox[id] = 1
		} else {
			splitBox[id] += 1
		}
	}
	return splitBox
}

func cubeBlogCollectRedisUpdate(splitBox map[string]int) {
	for k, v := range splitBox {
		redis.HSet("blog_profile_"+k, "collect", fmt.Sprintf("%v", v))
	}
}

func cubeBlogCollectDbUpdate(splitBox map[string]int) {
	b := new(blog.Blog)
	for k, v := range splitBox {
		b.Id, _ = strconv.Atoi(k)
		b.Collect = v
		_, err := database.Update(b, "collect")
		if err != nil {
			log.Error(err)
		}
	}
}
