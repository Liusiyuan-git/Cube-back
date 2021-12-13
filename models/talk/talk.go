package talk

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/user"
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

type Talk struct {
	Id      int
	CubeId  string `orm:"index"`
	Text    string `orm:"type(text)"`
	Date    string `orm:"index;type(datetime)"`
	Images  string `orm:"type(text)"`
	Love    int
	Comment int
}

type DeleteTalk struct {
	Id     int
	TalkId int
	Date   string
	CubeId string
	Images string `orm:"type(text)"`
}

func (b *Talk) TalkGet(mode, page string) (interface{}, interface{}, int64, string, bool) {
	var dataBlock []map[string]interface{}
	var countBox []string
	var countBlock []interface{}
	talkData, length := talkRedisGet(mode, page)
	if len(talkData) == 0 {
		talkDb, length, pass := talkDbGet(mode)
		return talkDb, countBlock, length, "db", pass
	}
	for _, item := range talkData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		id := fmt.Sprintf("%v", m["id"])
		countBox = append(countBox, id+"_like", id+"_comment")
		dataBlock = append(dataBlock, m)
	}
	countBlock = redis.HMGet("talk_like_and_comment", countBox)
	return dataBlock, countBlock, length, "redis", true
}

func (b *Talk) TalkSend(cubeid, text, images string) (string, bool) {
	talkImages, msg, pass := talkSendImages(cubeid, images)
	if !pass {
		return msg, pass
	}
	b.Id = 0
	b.CubeId = cubeid
	b.Comment = 0
	b.Text = text
	b.Love = 0
	b.Date = time.Now().Format("2006-01-02 15:04:05")
	b.Images = talkImages
	talkid, err := database.Insert(b)
	if err != nil {
		return "发送出错", false
	}
	talkSendRedis(talkid, *b)
	go talkMessageSend(talkid, b)
	return "", true
}

func (b *Talk) TalkDelete(images, cubeId, talkId, index, date string) (string, bool) {
	msg, pass := user.NumberCorrect(cubeId, talkId)
	if !pass {
		return msg, false
	}
	result, pass := talkDeleteDb(images, cubeId, talkId, date)
	if !pass {
		return result, pass
	}
	talkDeleteRedis(index, talkId, cubeId)
	return "", true
}

func talkMessageSend(talkid int64, b *Talk) {
	talkMessageDetailSet(talkid, b)
	talkMessageSendDb(talkid, b)
}

func talkSendImages(cubeid, images string) (string, string, bool) {
	var m []string
	var imagelist []string
	json.Unmarshal([]byte(images), &m)
	for index, image := range m {
		t, data, pass := base64Decode(image)
		bsonid := bson.NewObjectId()
		timeSplit := strings.Split(time.Now().Format("2006-01-02"), "-")
		timeJoin := strings.Join(timeSplit, "")
		filename := fmt.Sprintf("images%s%d.%s", bsonid.Hex(), index, t)
		filepath := fmt.Sprintf("/home/cube/images/talk/%s/%s", cubeid, timeJoin)
		pass = ssh.UploadFile(filename, filepath, data)
		if !pass {
			imagesRemove(imagelist)
			return "", "发送错误", false
		}
		imagelist = append(imagelist, filename)
	}
	return strings.Join(imagelist, ":"), "", true
}

func imagesRemove(images []string) {
	for _, item := range images {
		ssh.RemoveFile(item)
	}
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

func (b *Talk) TalkLike(talkid string) (string, bool) {
	redis.HIncrBy("talk_like_and_comment", talkid+"_like", 1)
	b.Id, _ = strconv.Atoi(talkid)
	b.Love, _ = strconv.Atoi(redis.HGet("talk_like_and_comment", talkid+"_like"))
	_, err := database.Update(b, "love")
	if err != nil {
		return "未知錯誤", false
	}

	return "", true
}

func (b *Talk) TalkSearch(keyWord, page string) (int, interface{}) {
	return talkEsSearch(keyWord, page)
}
