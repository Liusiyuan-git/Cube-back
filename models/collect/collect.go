package collect

import (
	"Cube-back/models/user"
	"Cube-back/redis"
	"encoding/json"
)

type Collect struct {
	Id     int
	CubeId string `orm:"index"`
	BlogId int    `orm:"index"`
}

func (o *Collect) BlogCollect(cubeid, blogid, cover, date, title, labelType string) (string, bool) {
	err := o.BlogCollectDb(cubeid, blogid)
	if err != nil {
		return "收藏错误", false
	}
	BlogCollectRedis(cubeid, blogid, cover, date, title, labelType)
	return "", true
}

func (o *Collect) BlogCollectConfirm(id, cubeid string) bool {
	_, pass := user.NumberCorrect(id, cubeid)
	if !pass {
		return false
	}
	ok := BlogCollectConfirmRedisGet(id, cubeid)
	if ok {
		return ok
	}
	ok = BlogCollectConfirmDbGet(id, cubeid)
	return ok
}

func (o *Collect) CollectProfileGet(blogId string) (interface{}, bool) {
	return redis.HMGet("blog_profile_"+blogId, []string{"love", "collect"}), true
}

func (o *Collect) CollectDelete(index, blogId, cubeId string) (string, bool) {
	msg, pass := user.NumberCorrect(blogId, cubeId)
	if !pass {
		return msg, pass
	}
	result, pass := collectDeleteDb(blogId, cubeId)
	if !pass {
		return result, pass
	}
	collectDeleteRedis(index, blogId, cubeId)
	return "", true
}

func (o *Collect) BlogCollectionGet(cubeid string) (interface{}, int64, bool) {
	_, pass := user.NumberCorrect(cubeid)
	if !pass {
		return []map[string]interface{}{}, 0, false
	}
	var dataBlock []map[string]interface{}
	collectData, length := collectRedisGet(cubeid)
	if len(collectData) == 0 {
		blogDb, length, pass := collectDbGet(cubeid)
		return blogDb, length, pass
	}
	for _, item := range collectData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, int64(length), true
}
