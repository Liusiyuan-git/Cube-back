package draft

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/ssh"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/siddontang/go/bson"
	"strconv"
	"strings"
)

type Draft struct {
	Id      int
	CubeId  string `orm:"index;pk"`
	Cover   string `orm:"type(text)"`
	Title   string
	Content string `orm:"type(text)"`
	Image   string `orm:"type(text)"`
}

func (b *Draft) DraftSend(cubeid, cover, title, content, images string) (string, bool) {
	sftpClient := ssh.RemoveDirectory("/home/cube/images/draft/" + cubeid)
	sftpClient.Wait()
	b.CubeId = cubeid
	id, draftExist := draftConfirm(cubeid)
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
	idi, err := strconv.Atoi(id)
	b.Id = idi
	b.Cover = coverName
	b.Title = title
	b.Content = content
	b.Image = contentImage
	if !draftExist {
		_, err := database.Insert(b)
		if err != nil {
			return "草稿保存错误", false
		}
		return "", true
	} else {
		_, err = database.Update(b, "cover", "title", "content", "image")
		if err != nil {
			return "草稿保存错误", false
		}
	}
	return "", true

}

func (b *Draft) DraftGet(cubeId string) (interface{}, bool) {
	cmd := "select * from draft where cube_id = ?"
	_, maps, pass := database.DBValues(cmd, cubeId)
	if !pass {
		return "", false
	} else {
		return maps, true
	}
}

func imageSave(cubeid, code string) (string, string, bool) {
	var box [][]string
	var imagelist []string
	json.Unmarshal([]byte(code), &box)
	for index, list := range box {
		for _, image := range list {
			t, data, pass := base64Decode(image)
			if !pass {
				return "", "草稿保存错误", false
			}
			bsonid := bson.NewObjectId()
			filename := fmt.Sprintf("content%s%d.%s", bsonid.Hex(), index, t)
			filepath := fmt.Sprintf("/home/cube/images/draft/%s", cubeid)
			pass = ssh.UploadFile(filename, filepath, data)
			if !pass {
				imagesRemove(imagelist)
				return "", "草稿保存错误", false
			}
			imagelist = append(imagelist, filename)
		}
	}
	return strings.Join(imagelist, ":"), "", true
}

func coverSave(cubeid, code string) (string, string, bool) {
	t, data, pass := base64Decode(code)
	if !pass {
		return "", "草稿保存错误", false
	}
	bsonid := bson.NewObjectId()
	filename := fmt.Sprintf("cover%s.%s", bsonid.Hex(), t)
	filepath := fmt.Sprintf("/home/cube/images/draft/%s", cubeid)
	pass = ssh.UploadFile(filename, filepath, data)
	if !pass {
		imagesRemove([]string{filepath + filename})
		return "", "草稿保存错误", false
	}
	return filename, "", true
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

func draftConfirm(cubeId string) (string, bool) {
	cmd := "select * from draft where cube_id = ?"
	num, maps, _ := database.DBValues(cmd, cubeId)
	if num > 0 && maps[0]["cube_id"] == cubeId {
		id := fmt.Sprintf("%v", maps[0]["id"])
		return id, true
	} else {
		return "", false
	}
}

func (b *Draft) DraftRemove(cubeId string) (interface{}, bool) {
	b.CubeId = cubeId
	//_, err := database.Delete(b, "id", "cube_id", "cover", "title", "content", "image")
	_, err := database.Delete(b)
	if err != nil {
		return "草稿删除错误", false
	}
	sftpClient := ssh.RemoveDirectory("/home/cube/images/draft/" + cubeId)
	sftpClient.Wait()
	return "", true
}
