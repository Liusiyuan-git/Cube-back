package blog

import (
	"Cube-back/database"
	"Cube-back/models/draft"
	"Cube-back/models/user"
	"Cube-back/redis"
	"encoding/json"
	"math"
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

type DeleteBlog struct {
	Id     int
	BlogId int
	CubeId string
	Image  string `orm:"type(text)"`
	Cover  string `orm:"type(text)"`
	Date   string `orm:"index;type(datetime)"`
}

func (b *Blog) BlogSend(cubeid, cover, title, content, text, images, label, labelType string) (string, bool) {
	msg, pass := user.NumberCorrect(cubeid)
	if !pass {
		return msg, false
	}
	var contentImage string
	contentImage = imageSave(images)
	b.Id = 0
	b.CubeId = cubeid
	b.Cover = cover
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
	go draftHandle(cubeid, cover, contentImage)
	go blogMessageSend(blogid, b)
	return "", true
}

func draftHandle(cubeid, cover, contentImage string) {
	r := new(draft.Draft)
	r.DraftImageMove(cubeid, cover+":"+contentImage)
	r.DraftRemove(cubeid)
}

func blogMessageSend(blogid int64, b *Blog) {
	blogMessageDetailSet(blogid, b)
	blogMessageSendDb(blogid, b)
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

func (b *Blog) BlogGet(mode, page, label, labeltype string) (interface{}, int64, string, bool) {
	blogData, length := blodRedisGet(mode, page, label, labeltype)
	if len(blogData) == 0 {
		blogDb, length, pass := blogDbGet(mode, label, labeltype)
		return blogDb, length, "db", pass
	}
	return blogData, length, "redis", true
}

func (b *Blog) BlogProfileGet(ids string) interface{} {
	return blogProfileRedisGet(ids)
}

func (b *Blog) BlogDelete(date, cover, image, label, labelType, index, blogId, cubeId string) (string, bool) {
	msg, pass := user.NumberCorrect(blogId, cubeId)
	if !pass {
		return msg, pass
	}
	result, pass := blogDeleteDb(date, cover, image, blogId, cubeId)
	if !pass {
		return result, pass
	}
	blogDeleteRedis(label, labelType, index, blogId, cubeId)
	return "", true
}

func (b *Blog) BlogForumGet(mode, page string) (interface{}, int64, bool) {
	var dataBlock []map[string]interface{}
	var t []string
	p, _ := strconv.ParseInt(page, 10, 64)
	if mode == "all" {
		mode = "new"
	}
	l := redis.LLen("blog_" + mode)
	if p > int64(math.Ceil(float64(l)/10)) {
		p = 1
	}
	t = redis.LRange("blog_"+mode, int64(0+(p-1)*10), int64(0+p*10-1))
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
	_, pass := user.NumberCorrect(id)
	if !pass {
		return "", false
	}
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

func imageSave(filenameBox string) string {
	var box [][]string
	var imagelist []string
	json.Unmarshal([]byte(filenameBox), &box)
	for _, list := range box {
		for _, filename := range list {
			imagelist = append(imagelist, filename)
		}
	}
	return strings.Join(imagelist, ":")
}

func (b *Blog) BlogSearch(keyWord, page string) (int, interface{}) {
	return blogEsSearch(keyWord, page)
}
