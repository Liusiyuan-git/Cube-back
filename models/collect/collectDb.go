package collect

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
	"strconv"
)

func (o *Collect) BlogCollectDb(cubeid, blogid string) error {
	id, _ := strconv.Atoi(blogid)
	o.Id = 0
	o.CubeId = cubeid
	o.BlogId = id
	_, err := database.Insert(o)
	if err != nil {
		return err
	}
	return err
}

func BlogCollectConfirmDbGet(id, cubeid string) bool {
	cmd := `select 1 from collect where blog_id = ? and cube_id = ? limit 1;`
	_, maps, pass := database.DBValues(cmd, id, cubeid)
	if !pass {
		return false
	} else if len(maps) == 1 {
		redis.HSet("blog_profile_"+id, cubeid, "1")
		return true
	} else {
		return false
	}
}

func collectDbGet(cubeid string) (interface{}, int64, bool) {
	key := "user_collect_" + cubeid
	cmd := `SELECT b.id, b.cube_id, b.title, b.cover, b.date, b.title, b.label_type FROM collect a INNER JOIN blog b ON a.blog_id = b.id  AND a.cube_id = ? ORDER BY a.id DESC`
	num, maps, pass := database.DBValues(cmd, cubeid)
	if num != 0 && !pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		if len(maps) >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return maps, num, pass
}
