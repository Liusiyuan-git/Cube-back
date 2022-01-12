package gron

import (
	"Cube-back/database"
	"Cube-back/elasticsearch"
	"Cube-back/log"
	"Cube-back/models/blog"
	"Cube-back/redis"
	"Cube-back/ssh"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
	"strings"
)

var label = map[string][]string{
	"language":       []string{"python", "go", "java", "javaScript", "c++", "c"},
	"middleware":     []string{"redis", "rabbitmq"},
	"virtualization": []string{"docker", "kubernetes", "microServices"},
	"database":       []string{"mysql"},
	"basics":         []string{"network", "dataStructure", "operatingSystem", "computerComposition"},
	"other":          []string{"live"},
}

var typeLibrary = map[string]string{
	"python":              "Python",
	"go":                  "Go",
	"java":                "Java",
	"javaScript":          "JavaScript",
	"c++":                 "C++",
	"c":                   "C",
	"redis":               "Redis",
	"rabbitmq":            "Rabbitmq",
	"docker":              "Docker",
	"kubernetes":          "kubernetes",
	"microServices":       "微服务",
	"mysql":               "Mysql",
	"network":             "网络",
	"dataStructure":       "数据结构和算法",
	"operatingSystem":     "操作系统",
	"computerComposition": "计算机组成原理",
	"live":                "生活",
}

func cubeBlogNewUpdate() {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name, b.image as user_image FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	num, maps, pass := database.DBValues(cmd)
	if pass {
		labelType := labelTypeBuild()
		cubeBlogDataSplit(maps, labelType)
		cubeBlogDataSet(num, labelType["all"], "new")
		cubeBlogDataFilterSet("new", labelType)
		cubeBlogProfileRedisSet(maps)
		cubeBlogEsSet(int(num), maps)
	}
}

func labelTypeBuild() map[string][]string {
	return map[string][]string{
		"all":                 []string{},
		"python":              []string{},
		"go":                  []string{},
		"java":                []string{},
		"javaScript":          []string{},
		"c++":                 []string{},
		"c":                   []string{},
		"redis":               []string{},
		"rabbitmq":            []string{},
		"docker":              []string{},
		"kubernetes":          []string{},
		"mysql":               []string{},
		"microServices":       []string{},
		"network":             []string{},
		"dataStructure":       []string{},
		"operatingSystem":     []string{},
		"computerComposition": []string{},
		"live":                []string{},
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
		go cubeBlogCommentUpdate(fmt.Sprintf("%v", item["id"]))
	}
}

func cubeBlogDetailUpdate() {
	txpipeline := redis.TxPipeline()
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.image, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name, b.image as user_image FROM blog a inner join user b on a.cube_id = b.cube_id order by a.id desc`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		for _, item := range maps {
			id := fmt.Sprintf("%v", item["id"])
			bjson, _ := json.Marshal(item)
			txpipeline.HSet("blog_detail", id, string(bjson))
		}
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func cubeBlogCommentUpdate(blogid string) {
	cmd := `SELECT a.id, a.cube_id, a.comment, a.date, a.love, b.image, b.name FROM blog_comment a INNER JOIN user b ON a.cube_id = b.cube_id WHERE a.blog_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, blogid)
	key := "blog_detail_comment_" + blogid
	if pass {
		txpipeline := redis.TxPipeline()
		if num != 0 {
			var l = redis.LLen(key)
			if l <= num {
				for i := int64(0); i < num; i++ {
					bjson, _ := json.Marshal(maps[i])
					redisValue := string(bjson)
					if i+1 > l {
						txpipeline.RPush(key, redisValue)
					} else {
						txpipeline.LSet(key, i, redisValue)
					}
				}
			} else {
				for i := int64(0); i < num; i++ {
					bjson, _ := json.Marshal(maps[i])
					redisValue := string(bjson)
					txpipeline.LSet(key, i, redisValue)
				}
				txpipeline.LTrim(key, 0, num-1)
			}
		} else {
			txpipeline.LTrim(key, 1, 0)
		}
		txpipeline.HSet("blog_profile_"+blogid, "comment", fmt.Sprintf("%v", num))
		txpipeline.Exec()
		txpipeline.Close()
		cubeBlogProfileCommentDbSet(blogid, num)
	}
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
	txpipeline := redis.TxPipeline()
	if l <= num {
		for i := int64(0); i < num; i++ {
			if i+1 > l {
				txpipeline.RPush(key, maps[i])
			} else {
				txpipeline.LSet(key, i, maps[i])
			}
		}
	} else {
		for i := int64(0); i < num; i++ {
			txpipeline.LSet(key, i, maps[i])
		}
		txpipeline.LTrim(key, 0, num-1)
	}
	txpipeline.Exec()
	txpipeline.Close()
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

func cubeBlogProfileRedisSet(maps []orm.Params) {
	txpipeline := redis.TxPipeline()
	for _, item := range maps {
		id := fmt.Sprintf("%v", item["id"])
		txpipeline.HMSet("blog_profile_"+id, map[string]interface{}{"love": item["love"], "view": item["view"]})
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func cubeBlogEsSet(num int, maps []orm.Params) {
	EsLen, EsMaps := elasticsearch.Client.SearchAll("blog")
	if num >= EsLen {
		for index, item := range maps {
			var box = map[string]interface{}{}
			box["label_type"] = typeLibrary[item["label_type"].(string)]
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
			elasticsearch.Client.Create("blog", redisValue, index)
		}
	} else {
		for index, item := range EsMaps {
			if (index + 1) <= num {
				var box = map[string]interface{}{}
				box["label_type"] = typeLibrary[maps[index]["label_type"].(string)]
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
				elasticsearch.Client.Create("blog", redisValue, index)
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
	if pass {
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
	txpipeline := redis.TxPipeline()
	for k, v := range splitBox {
		txpipeline.HSet("blog_profile_"+k, "collect", fmt.Sprintf("%v", v))
	}
	txpipeline.Exec()
	txpipeline.Close()
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

func cubeBlogCleanAll() {
	cmd := `select * from delete_blog`
	_, maps, pass := database.DBValues(cmd)
	if pass {
		cubeBlogCleanRedisAll(maps)
		cubeBlogCleanImageAll(maps)
		cubeBlogCleanConfirm()
	}
}

func cubeBlogCleanConfirm() {
	cmd := `truncate table delete_blog`
	database.DBValues(cmd)
}

func cubeBlogCleanRedisAll(maps []orm.Params) {
	txpipeline := redis.TxPipeline()
	for _, item := range maps {
		blogId, _ := item["blog_id"].(string)
		txpipeline.HDel("blog_detail", blogId)
		txpipeline.HDel("blog_message_detail", "cover_"+blogId)
		txpipeline.HDel("blog_message_detail", "title_"+blogId)
		txpipeline.HDel("blog_message_detail", "date_"+blogId)
		txpipeline.HDel("blog_message_detail", "type_"+blogId)
		txpipeline.Del("blog_detail_comment_" + blogId)
		txpipeline.Del("blog_profile_" + blogId)
		txpipeline.Del("blog_detail_" + blogId + "_get")
		txpipeline.Del("blog_detail_" + blogId + "_comment_get")
	}
	txpipeline.Exec()
	txpipeline.Close()
}

func cubeBlogCleanImageAll(maps []orm.Params) {
	var deleteBox []string
	var deleteString string
	for _, item := range maps {
		cubeId, _ := item["cube_id"].(string)
		date := strings.Join(strings.Split(strings.Split(item["date"].(string), " ")[0], "-"), "")
		cover := item["cover"].(string)
		image := item["image"].(string)
		imagePath := "/home/cube/images/blog/" + cubeId + "/" + date
		for _, each := range strings.Split(cover+":"+image, ":") {
			if each != "" {
				deleteBox = append(deleteBox, each)
			}
		}
		deleteString = strings.Join(deleteBox, " ")
		if deleteString != "" {
			ssh.CommandExecute("cd " + imagePath + ";" + "rm -rf " + deleteString)
		}
	}
}
