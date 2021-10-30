package blog

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/draft"
	"Cube-back/redis"
	"Cube-back/ssh"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/siddontang/go/bson"
	"strconv"
	"strings"
	"time"
)

type Blog struct {
	Id        int
	CubeId    string `orm:"index"`
	Cover     string `orm:"type(text)"`
	Title     string
	Content   string `orm:"type(text)"`
	Text      string `orm:"type(text)"`
	Image     string `orm:"type(text)"`
	Date      string `orm:"index;type(datetime)"`
	Label     string
	LabelType string
	Love      int
	Comment   int
	Collect   int
	View      int
}

func (b *Blog) BlogSend(cubeid, cover, title, content, text, images, label, labelType string) (string, bool) {
	var coverName string
	var contentImage string
	var msg string
	var pass bool
	if cover != "" {
		coverName, msg, pass = coverSave(cubeid, cover)
		if !pass {
			return msg, false
		}
	}
	contentImage, msg, pass = imageSave(cubeid, images)
	if !pass {
		return msg, false
	}
	b.Id = 0
	b.CubeId = cubeid
	b.Cover = coverName
	b.Title = title
	b.Content = content
	b.Text = text
	b.Image = contentImage
	b.Date = blogTime()
	b.Label = label
	b.LabelType = labelType
	b.Love = 0
	b.Comment = 0
	b.Collect = 0
	b.View = 0
	blogid, err := database.Insert(b)
	if err != nil {
		return "发送错误", false
	}
	blogSendRedis(blogid, *b)
	r := new(draft.Draft)
	r.DraftRemove(cubeid)
	go blogMessageSend(b)
	return "", true
}

func blogMessageSend(b *Blog) {
	blogMessageSendDb(b)
	blogMessageSendRedis(b)
}

func blogTime() string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t
}

func (b *Blog) BlogLike(id string) (string, bool) {
	key := "blog_profile_" + id
	redis.HIncrBy(key, "love", 1)
	b.Id, _ = strconv.Atoi(id)
	b.Love, _ = strconv.Atoi(redis.HGet(key, "love"))
	_, err := database.Update(b, "love")
	if err != nil {
		redis.HIncrBy(key, "love", -1)
		return "未知错误", false
	}
	return "", true
}

func (b *Blog) BlogView(id string) (string, bool) {
	key := "blog_profile_" + id
	redis.HIncrBy(key, "view", 1)
	b.Id, _ = strconv.Atoi(id)
	b.View, _ = strconv.Atoi(redis.HGet(key, "view"))
	_, err := database.Update(b, "view")
	if err != nil {
		redis.HIncrBy(key, "view", -1)
		return "未知错误", false
	}
	return "", true
}

func (b *Blog) BlogGet(mode, page, label, labeltype string) (interface{}, interface{}, int64, string, bool) {
	var dataBlock []map[string]interface{}
	var countBox [][]interface{}
	blogData, length := blodRedisGet(mode, page, label, labeltype)
	if len(blogData) == 0 {
		blogDb, length, pass := blogDbGet(mode, label, labeltype)
		return blogDb, countBox, length, "db", pass
	}
	for _, item := range blogData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		id := fmt.Sprintf("%v", m["id"])
		countBox = append(countBox, redis.HMGet("blog_profile_"+id, []string{"love", "comment", "collect", "view"}))
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, countBox, length, "redis", true
}

func (b *Blog) BlogForumGet(mode, page string) (interface{}, int64, bool) {
	var dataBlock []map[string]interface{}
	var t []string
	p, _ := strconv.Atoi(page)
	if mode == "all" {
		mode = "new"
	}
	t = redis.LRange("blog_"+mode, int64(0+(p-1)*10), int64(0+p*10-1))
	l := redis.LLen("blog_" + mode)
	if t == nil {
		return t, 0, false
	}
	for _, item := range t {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, l, true
}

func (b *Blog) BlogDetail(id string) (interface{}, bool) {
	key := "blog_detail_" + id + "_get"
	result, pass := blogDetailRedisGet(id)
	if pass {
		return result, true
	}
	if "true" == blogRedisLockStatus(key) {
		return "", false
	}
	blogRedisLock(key, "true")
	result, pass = blogDetailDbGet(id, b)
	blogRedisLock(key, "false")
	if pass {
		return result, true
	}
	return "", false
}

func coverSave(cubeid, code string) (string, string, bool) {
	t, data, pass := base64Decode(code)
	if !pass {
		return "", "发送错误", false
	}
	bsonid := bson.NewObjectId()
	timeSplit := strings.Split(time.Now().Format("2006-01-02"), "-")
	timeJoin := strings.Join(timeSplit, "")
	filename := fmt.Sprintf("cover%s.%s", bsonid.Hex(), t)
	filepath := fmt.Sprintf("/home/cube/images/blog/%s/%s", cubeid, timeJoin)
	pass = ssh.UploadFile(filename, filepath, data)
	if !pass {
		imagesRemove([]string{filepath + filename})
		return "", "发送错误", false
	}
	return filename, "", true
}

func imageSave(cubeid, code string) (string, string, bool) {
	var box [][]string
	var imagelist []string
	json.Unmarshal([]byte(code), &box)
	for index, list := range box {
		for _, image := range list {
			t, data, pass := base64Decode(image)
			if !pass {
				return "", "发送错误", false
			}
			bsonid := bson.NewObjectId()
			timeSplit := strings.Split(time.Now().Format("2006-01-02"), "-")
			timeJoin := strings.Join(timeSplit, "")
			filename := fmt.Sprintf("content%s%d.%s", bsonid.Hex(), index, t)
			filepath := fmt.Sprintf("/home/cube/images/blog/%s/%s", cubeid, timeJoin)
			pass = ssh.UploadFile(filename, filepath, data)
			if !pass {
				imagesRemove(imagelist)
				return "", "发送错误", false
			}
			imagelist = append(imagelist, filename)
		}
	}
	return strings.Join(imagelist, ":"), "", true
}

func base64Decode(code string) (string, []uint8, bool) {
	s := strings.Split(code, "data:image/")
	t := strings.Split(s[1], ";")
	enc := base64.StdEncoding
	data, err := enc.DecodeString(t[1][7:])
	if err != nil {
		log.Error(err)
		return "", make([]uint8, 1), false
	} else {
		return t[0], data, true
	}
}

func imagesRemove(images []string) {
	for _, item := range images {
		ssh.RemoveFile(item)
	}
}
