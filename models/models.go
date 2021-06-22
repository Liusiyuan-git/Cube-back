package models

import (
	"Cube-back/models/blog"
	"Cube-back/models/draft"
	"Cube-back/models/user"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func ModelRegister() {
	orm.RegisterModel(new(user.User))
	orm.RegisterModel(new(blog.Blog))
	orm.RegisterModel(new(draft.Draft))
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
