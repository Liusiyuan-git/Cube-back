package models

import (
	"Cube-back/models/blog"
	"Cube-back/models/blogcomment"
	"Cube-back/models/care"
	"Cube-back/models/collect"
	"Cube-back/models/draft"
	"Cube-back/models/talk"
	"Cube-back/models/talkcomment"
	"Cube-back/models/user"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func ModelRegister() {
	orm.RegisterModel(new(user.User))
	orm.RegisterModel(new(blog.Blog))
	orm.RegisterModel(new(draft.Draft))
	orm.RegisterModel(new(collect.Collect))
	orm.RegisterModel(new(blogcomment.BlogComment))
	orm.RegisterModel(new(talk.Talk))
	orm.RegisterModel(new(talkcomment.TalkComment))
	orm.RegisterModel(new(care.Care))
	orm.RunCommand()
}

func RunSyncdb() {
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		logs.Error(err)
	}
}

func init() {
	ModelRegister()
	RunSyncdb()
}
