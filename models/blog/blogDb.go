package blog

import (
	"Cube-back/database"
	"Cube-back/redis"
	"encoding/json"
)

func blogDbGet(mode, label, labeltype string) (interface{}, int64, bool) {
	var key = blogDbKeySet(mode, label, labeltype)
	if "true" == blogRedisLockStatus(key+"_get") {
		return "数据更新中，请稍后再试", 0, false
	}
	blogRedisLock(key+"_get", "true")
	var cmd = blogDbCmdFilterSet(label, labeltype)
	cmd = blogDbCmdModeSet(cmd, mode)
	num, maps, pass := database.DBValues(cmd)
	blogRedisLock(key+"_get", "false")
	if num != 0 && pass {
		for _, item := range maps {
			bjson, _ := json.Marshal(item)
			redisValue := string(bjson)
			redis.RPush(key, redisValue)
		}
		redis.Set(key+"_get", "false")
		if len(maps) >= 10 {
			return maps[0:9], num, true
		} else {
			return maps[0:], num, true
		}
	}
	return "", 0, false
}

func blogDbKeySet(mode, label, labeltype string) string {
	var key string
	if labeltype == "" {
		key = "blog_" + mode
	} else if labeltype == "all" {
		key = "blog_" + label + "_all_" + mode
	} else {
		key = "blog_" + labeltype + "_" + mode
	}
	return key
}

func blogDbCmdFilterSet(label, labeltype string) string {
	var cmd = `select a.id, a.cube_id, a.cover, a.title, a.text, a.date, a.label, a.label_type, a.love, a.comment, a.collect,
	a.view, b.name FROM blog a inner join user b on a.cube_id = b.cube_id`
	if labeltype != "" {
		switch labeltype {
		case "all":
			cmd += " and a.label = '" + label + "'"
		default:
			cmd += " and a.label_type = '" + labeltype + "'"
		}
	}
	return cmd
}

func blogDbCmdModeSet(cmd, mode string) string {
	switch mode {
	case "new":
		cmd += " order by id desc"
	case "hot":
		cmd += " order by a.love desc"
	}
	return cmd
}

func blogDetailDbGet(id string, b *Blog) (interface{}, bool) {
	cmd := `select a.id, a.cube_id, a.cover, a.title, a.content, a.image, a.date, a.love, a.collect, a.view, a.comment ,b.name from blog a inner join user b on a.cube_id = b.cube_id and a.id = ? order by a.id desc`
	_, maps, pass := database.DBValues(cmd, id)
	if !pass {
		return "", false
	} else {
		if maps != nil {
			bjson, _ := json.Marshal(maps[0])
			redisValue := string(bjson)
			redis.HSet("blog_detail", id, redisValue)
			return maps, true
		} else {
			return "", false
		}
	}
}
