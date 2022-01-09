package draft

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/user"
	"Cube-back/ssh"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/siddontang/go/bson"
	"strconv"
	"strings"
	"time"
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
	b.CubeId = cubeid
	id, draftExist := draftConfirm(cubeid)
	var contentImage string
	contentImage = imageSave(images)
	idi, err := strconv.Atoi(id)
	b.Id = idi
	b.Cover = cover
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
	_, pass := user.NumberCorrect(cubeId)
	if !pass {
		return "", false
	}
	cmd := "select * from draft where cube_id = ?"
	_, maps, pass := database.DBValues(cmd, cubeId)
	if !pass {
		return "", false
	} else {
		return maps, true
	}
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

func (b *Draft) DraftImageUpload(cubeId, code, mode string) (string, string, bool) {
	t, data, pass := base64Decode(code)
	if !pass {
		return "", "图片上传错误", false
	}
	bsonid := bson.NewObjectId()
	filename := fmt.Sprintf(mode+"%s.%s", bsonid.Hex(), t)
	filepath := fmt.Sprintf("/home/cube/images/draft/%s", cubeId)
	pass = ssh.UploadFile(filename, filepath, data)
	if !pass {
		imagesRemove([]string{filepath + filename})
		return "", "图片上传错误", false
	}
	ssh.CommandExecute("cd /home/lsy;mv text1 text2 text3 ../")
	return filename, "", true
}

func (b *Draft) DraftImageDelete(cubeId, filename string) bool {
	filepath := fmt.Sprintf("/home/cube/images/draft/%s/%s", cubeId, filename)
	ssh.RemoveFile(filepath)
	return true
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
	_, pass := user.NumberCorrect(cubeId)
	if !pass {
		return "", false
	}
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
	_, err := database.Delete(b)
	if err != nil {
		return "草稿删除错误", false
	}
	sftpClient := ssh.RemoveDirectory("/home/cube/images/draft/" + cubeId)
	sftpClient.Wait()
	return "", true
}

func (b *Draft) DraftImageMove(cubeId, filename string) {
	var eBox []string
	timeSplit := strings.Split(time.Now().Format("2006-01-02"), "-")
	timeJoin := strings.Join(timeSplit, "")
	for _, item := range strings.Split(filename, ":") {
		eBox = append(eBox, item)
	}
	directoryPath := "/home/cube/images/blog/" + cubeId + "/" + timeJoin
	draftPath := "/home/cube/images/draft/" + cubeId
	ssh.CommandExecute("mkdir -p" + directoryPath)
	ssh.CommandExecute("cd " + draftPath + ";" + "mv " + strings.Join(eBox, " ") + " " + directoryPath)
}
