package collect

import (
	"Cube-back/database"
	"Cube-back/models/blog"
	"strconv"
)

func (o *Collect) BlogCollectDb(cudeid, blogid, collect string) error {
	b := new(blog.Blog)
	id, _ := strconv.Atoi(blogid)
	c, _ := strconv.Atoi(collect)
	o.Id = 0
	o.CubeId = cudeid
	o.BlogId = id
	b.Id = id
	b.Collect = c
	_, err := database.Insert(o)
	database.Update(b, "collect")
	return err
}

func BlogCollectConfirmDbGet(id, cubeid string) bool {
	cmd := `select 1 from collect where blog_id = ? and cube_id = ? limit 1;`
	_, maps, pass := database.DBValues(cmd, id, cubeid)
	if !pass {
		return false
	} else if len(maps) == 1 {
		blogCollectIdSet(cubeid, id)
		return true
	} else {
		return false
	}
}
