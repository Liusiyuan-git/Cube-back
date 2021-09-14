package common

import (
	"Cube-back/controllers"
	"Cube-back/models/blog"
	"Cube-back/models/blogcomment"
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

func (c *Controller) BlogGet() {
	data := c.RequestBodyData()
	mode := data["mode"]
	page := data["page"]
	label := data["label"]
	labelType := data["label_type"]
	content, length, pass := b.BlogGet(mode, page, label, labelType)
	result := make(map[string]interface{})
	result["content"] = content
	result["length"] = length
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
	content, length, pass := t.TalkGet(mode, page)
	result := make(map[string]interface{})
	result["content"] = content
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
	cubeid := data["id"]
	like := data["like"]
	content, pass := b.BlogLike(cubeid, like)
	result := make(map[string]interface{})
	result["msg"] = content
	result["like"] = like
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
	like := data["like"]
	index := data["index"]
	mode := data["mode"]
	content, pass := t.TalkLike(talkid, like, index, mode)
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
