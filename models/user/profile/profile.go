package profile

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"Cube-back/ssh"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/siddontang/go/bson"
	"strconv"
	"strings"
)

type Profile struct {
	user.User
}

func (p *Profile) PasswordChange(phone, password, code string) (string, bool) {
	s, pass, id := passwordParamsCheck(phone, password, code)
	if !pass {
		return s, pass
	}
	idi, err := strconv.Atoi(id)
	password = crypt.Set(password)
	u := new(user.User)
	u.Id, u.Password = idi, password
	_, err = database.Update(u, "password")
	if err != nil {
		fmt.Println(err)
		return "密码修改失败", false
	}
	return "", true
}

func (p *Profile) ProfileBlogGet(cubeId, page string) (interface{}, int64, bool) {
	var dataBlock []map[string]interface{}
	profileBlogData, length := profileBlogRedisGet(cubeId, page)
	if len(profileBlogData) == 0 {
		blogDb, length, pass := profileBlogDbGet(cubeId)
		return blogDb, length, pass
	}
	for _, item := range profileBlogData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, length, true
}

func (p *Profile) ProfileTalkGet(cubeId, page string) (interface{}, int64, bool) {
	var dataBlock []map[string]interface{}
	profileTalkData, length := profileTalkRedisGet(cubeId, page)
	if len(profileTalkData) == 0 {
		blogDb, length, pass := profileTalkDbGet(cubeId)
		return blogDb, length, pass
	}
	for _, item := range profileTalkData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return dataBlock, length, true
}

func passwordParamsCheck(phone, password, code string) (string, bool, string) {
	var pass bool
	var msg string
	msg, pass = paramsEmpty(password, phone, code)
	if !pass {
		return msg, pass, ""
	}
	msg, pass = user.CodeCorrect(phone, code)
	if !pass {
		return msg, pass, ""
	}
	msg, pass = user.PhoneConfirm(phone)
	if !pass {
		return msg, pass, ""
	}
	return "", true, msg
}

func paramsEmpty(password, phone, code string) (string, bool) {
	if password == "" || phone == "" || code == "" {
		return "表单信息不完整，请检查", false
	}
	return "", true
}

func (p *Profile) SendUserImage(cubeid, image string) (string, bool) {
	imageName, msg, pass := imageSave(cubeid, image)
	if !pass {
		return msg, false
	}
	u := new(user.User)
	u.Image = imageName
	u.CubeId = cubeid
	_, err := database.Update(u, "image")
	if err != nil {
		log.Error(err)
		return "未知错误", false
	}
	SendUserImageRedis(cubeid, imageName)
	return "", true
}

func (p *Profile) UserProfileGet(cubeId string) (interface{}, bool) {
	profileData := userProfileRedisGet(cubeId)
	if profileData == "" {
		userProfileDb, pass := userProfileBlogDbGet(cubeId)
		return userProfileDb, pass
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(profileData), &m)
	return m, true
}

func (p *Profile) UserIntroduceSend(cubeId, introduce string) bool {
	pass := UserIntroduceDbSend(cubeId, introduce)
	if !pass {
		return false
	}
	UserIntroduceRedisSend(cubeId, introduce)
	return true
}

func imageSave(cubeid, image string) (string, string, bool) {
	t, data, pass := base64Decode(image)
	if !pass {
		return "", "发送错误", false
	}
	bsonid := bson.NewObjectId()
	filename := fmt.Sprintf("userimage%s.%s", bsonid.Hex(), t)
	filepath := fmt.Sprintf("/home/cube/images/user/image/%s", cubeid)
	removeDirectory(filepath)
	pass = ssh.UploadFile(filename, filepath, data)
	if !pass {
		imagesRemove([]string{filepath + filename})
		return "", "发送错误", false
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

func removeDirectory(path string) {
	ssh.RemoveDirectory(path)
}
