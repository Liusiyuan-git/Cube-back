package collect

import (
	"Cube-back/database"
)

type Collect struct {
	Id     int
	CubeId string `orm:"index"`
	BlogId int    `orm:"index"`
}

func (o *Collect) BlogCollect(cudeid, blogid, collect string) (string, bool) {
	err := o.BlogCollectDb(cudeid, blogid, collect)
	if err != nil {
		return "收藏错误", false
	}
	BlogCollectRedis(cudeid, blogid, collect)
	return "", true
}

func (o *Collect) BlogCollectConfirm(id, cubeid string) bool {
	ok := BlogCollectConfirmRedisGet(id, cubeid)
	if ok {
		return ok
	}
	ok = BlogCollectConfirmDbGet(id, cubeid)
	return ok
}

func (o *Collect) BlogCollectionGet(cubeid string) (interface{}, bool) {
	cmd := `SELECT a.blog_id, b.title FROM collect a INNER JOIN blog b ON a.blog_id = b.id  AND a.cube_id = ? ORDER BY a.id DESC`
	_, maps, pass := database.DBValues(cmd, cubeid)
	if !pass {
		return "", false
	} else {
		return maps, true
	}
}
