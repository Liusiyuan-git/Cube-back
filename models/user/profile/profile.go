package profile

import (
	"Cube-back/database"
	"Cube-back/log"
	"Cube-back/models/care"
	"Cube-back/models/common/crypt"
	"Cube-back/models/user"
	"Cube-back/redis"
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

func (p *Profile) ProfileTalkGet(cubeId, page string) (interface{}, interface{}, int64, string, bool) {
	var dataBlock []map[string]interface{}
	var countBox []string
	var countBlock []interface{}
	profileTalkData, length := profileTalkRedisGet(cubeId, page)
	if len(profileTalkData) == 0 {
		talkDb, length, pass := profileTalkDbGet(cubeId)
		return talkDb, countBlock, length, "db", pass
	}
	for _, item := range profileTalkData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		id := fmt.Sprintf("%v", m["id"])
		countBox = append(countBox, id+"_like", id+"_comment")
		dataBlock = append(dataBlock, m)
	}
	countBlock = redis.HMGet("talk_like_and_comment", countBox)
	return dataBlock, countBlock, length, "redis", true
}

func (p *Profile) ProfileCollectGet(cubeId, page string) (interface{}, int64, bool) {
	var dataBlock []map[string]interface{}
	collectData, length := profileCollectRedisGet(cubeId, page)
	if len(collectData) == 0 {
		blogDb, length, pass := profileCollectDbGet(cubeId)
		return blogDb, length, pass
	}
	for _, item := range collectData {
		var m map[string]interface{}
		json.Unmarshal([]byte(item), &m)
		dataBlock = append(dataBlock, m)
	}
	return collectData, length, true
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
		userProfileDb, pass := userProfileDbGet(cubeId)
		return userProfileDb, pass
	}
	return profileData, true
}

func (p *Profile) UserImageUpdate(cubeId string) (interface{}, bool) {
	profileData := userImageRedisGet(cubeId)
	return profileData, true
}

func (p *Profile) UserCareSet(id, cubeId string) (string, bool) {
	c := new(care.Care)
	c.Care = id
	c.Cared = cubeId
	_, err := database.Insert(c)
	if err != nil {
		log.Error(err)
		return "未知错误", true
	}
	userCareRedisSet(id, cubeId)
	userCareDbSet(id, cubeId)
	return "", true
}

func (p *Profile) UserCareGet(id, cubeId string) (interface{}, bool) {
	careRedisData := userCareRedisGet(id)
	var careBox = []map[string]string{}
	if len(careRedisData) == 0 {
		userCareData, pass := userCareDbGet(id)
		return userCareData, pass
	}
	for _, item := range careRedisData {
		profile := redis.HMGet("user_profile_"+item, []string{"image", "name", "introduce"})
		careBox = append(careBox, map[string]string{"cube_id": item, "image": fmt.Sprintf("%v", profile[0]), "name": fmt.Sprintf("%v", profile[1]), "introduce": fmt.Sprintf("%v", profile[2])})
	}
	return careBox, true
}

func (p *Profile) UserCareConfirm(id, cubeId string) (bool, bool) {
	result := redis.HExists("user_care_"+id, cubeId)
	if result {
		return true, true
	}
	cmd := `select * from care where care=? and cared=?`
	num, _, pass := database.DBValues(cmd, id, cubeId)
	if pass {
		if num != 0 {
			return true, true
		} else {
			return false, true
		}
	} else {
		return false, false
	}
}

func (p *Profile) UserCareCancel(id, cubeId string) bool {
	cmd := "DELETE FROM care where care=? and cared=?"
	_, _, pass := database.DBValues(cmd, id, cubeId)
	if !pass {
		return false
	} else {
		userCareRedisCancelSet(id, cubeId)
		userCareDbSet(id, cubeId)
		return true
	}
}

func (p *Profile) UserIntroduceSend(cubeId, introduce string) bool {
	pass := UserIntroduceDbSend(cubeId, introduce)
	if !pass {
		return false
	}
	UserIntroduceRedisSend(cubeId, introduce)
	return true
}

func (p *Profile) UserNameSend(cubeId, name string) bool {
	pass := UserNameDbSend(cubeId, name)
	if !pass {
		return false
	}
	UserNameRedisSend(cubeId, name)
	return true
}

func (p *Profile) UserProfileCare(cubeId string) (interface{}, bool) {
	profileCare := profileCareRedisGet(cubeId)
	if len(profileCare) == 0 {
		profileCare, pass := profileCareDbGet(cubeId)
		return profileCare, pass
	}
	return profileCare, true
}

func (p *Profile) UserProfileCared(cubeId string) (interface{}, bool) {
	profileCared := profileCaredRedisGet(cubeId)
	if len(profileCared) == 0 {
		profileCared, pass := profileCaredDbGet(cubeId)
		return profileCared, pass
	}
	return profileCared, true
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
