package common

import (
	"Cube-back/controllers"
	"Cube-back/models/blog"
	"Cube-back/models/blogcomment"
	"Cube-back/models/collect"
	"Cube-back/models/leaveMessage"
	"Cube-back/models/talk"
	"Cube-back/models/talkcomment"
	"Cube-back/models/user/profile"
)

type Controller struct {
	controllers.MainController
}

var b = new(blog.Blog)
var bc = new(blogcomment.BlogComment)
var t = new(talk.Talk)
var tc = new(talkcomment.TalkComment)
var p = new(profile.Profile)
var co = new(collect.Collect)
var l = new(leaveMessage.LeaveMessage)

func (c *Controller) BlogGet() {
	data := c.RequestBodyData()
	mode := data["mode"]
	page := data["page"]
	label := data["label"]
	labelType := data["label_type"]
	content, profile, length, mode, pass := b.BlogGet(mode, page, label, labelType)
	result := make(map[string]interface{})
	result["content"] = content
	result["profile"] = profile
	result["length"] = length
	result["mode"] = mode
	c.DataCallBack(result, pass)
}

func (c *Controller) BlogForumGet() {
	data := c.RequestBodyData()
	mode := data["mode"]
	page := data["page"]
	content, length, pass := b.BlogForumGet(mode, page)
	result := make(map[string]interface{})
	result["content"] = content
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) TalkGet() {
	data := c.RequestBodyData()
	mode := data["mode"]
	page := data["page"]
	content, count, length, mode, pass := t.TalkGet(mode, page)
	result := make(map[string]interface{})
	result["content"] = content
	result["count"] = count
	result["mode"] = mode
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) BlogDetail() {
	data := c.RequestBodyData()
	id := data["id"]
	result := make(map[string]interface{})
	if id == "" {
		result["msg"] = "未知错误"
		c.DataCallBack(result, false)
	} else {
		content, pass := b.BlogDetail(id)
		result["content"] = content
		c.DataCallBack(result, pass)
	}
}

func (c *Controller) BlogLike() {
	data := c.RequestBodyData()
	blogid := data["id"]
	content, pass := b.BlogLike(blogid)
	result := make(map[string]interface{})
	result["msg"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) BlogView() {
	data := c.RequestBodyData()
	blogid := data["id"]
	content, pass := b.BlogView(blogid)
	result := make(map[string]interface{})
	result["msg"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) BlogCommonLike() {
	data := c.RequestBodyData()
	commentid := data["id"]
	blogid := data["blogid"]
	like := data["like"]
	index := data["index"]
	content, pass := bc.BlogCommonLike(commentid, blogid, index, like)
	result := make(map[string]interface{})
	result["msg"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) TalkLike() {
	data := c.RequestBodyData()
	talkid := data["id"]
	content, pass := t.TalkLike(talkid)
	result := make(map[string]interface{})
	result["msg"] = content
	c.DataCallBack(result, pass)
}

func (c *Controller) BlogCommonGet() {
	data := c.RequestBodyData()
	blogid := data["id"]
	page := data["page"]
	content, length, pass := bc.BlogCommonGet(blogid, page)
	result := make(map[string]interface{})
	result["comment"] = content
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) TalkCommentGet() {
	data := c.RequestBodyData()
	talkid := data["id"]
	page := data["page"]
	content, length, pass := tc.TalkCommonGet(talkid, page)
	result := make(map[string]interface{})
	result["content"] = content
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) ProfileBlogGet() {
	data := c.RequestBodyData()
	cubeid := data["cube_id"]
	page := data["page"]
	content, length, pass := p.ProfileBlogGet(cubeid, page)
	result := make(map[string]interface{})
	result["content"] = content
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) ProfileTalkGet() {
	data := c.RequestBodyData()
	cubeid := data["cube_id"]
	page := data["page"]
	content, count, length, mode, pass := p.ProfileTalkGet(cubeid, page)
	result := make(map[string]interface{})
	result["content"] = content
	result["count"] = count
	result["mode"] = mode
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) ProfileCollectGet() {
	data := c.RequestBodyData()
	cubeid := data["cube_id"]
	page := data["page"]
	content, length, pass := p.ProfileCollectGet(cubeid, page)
	result := make(map[string]interface{})
	result["content"] = content
	result["length"] = length
	c.DataCallBack(result, pass)
}

func (c *Controller) UserProfileGet() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	profile, pass := p.UserProfileGet(cubeId)
	result := make(map[string]interface{})
	result["profile"] = profile
	c.DataCallBack(result, pass)
}

func (c *Controller) CollectProfileGet() {
	data := c.RequestBodyData()
	blogId := data["id"]
	profile, pass := co.CollectProfileGet(blogId)
	result := make(map[string]interface{})
	result["profile"] = profile
	c.DataCallBack(result, pass)
}

func (c *Controller) UserProfileCare() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	profileCare, pass := p.UserProfileCare(cubeId)
	result := make(map[string]interface{})
	result["profileCare"] = profileCare
	c.DataCallBack(result, pass)
}

func (c *Controller) UserProfileCared() {
	data := c.RequestBodyData()
	cubeId := data["cubeid"]
	profileCared, pass := p.UserProfileCared(cubeId)
	result := make(map[string]interface{})
	result["profileCared"] = profileCared
	c.DataCallBack(result, pass)
}

func (c *Controller) ProfileLeaveGet() {
	data := c.RequestBodyData()
	cubeId := data["cube_id"]
	page := data["page"]
	leaveData, length, pass := l.LeaveGet(cubeId, page)
	result := make(map[string]interface{})
	result["content"] = leaveData
	result["length"] = length
	c.DataCallBack(result, pass)
}
