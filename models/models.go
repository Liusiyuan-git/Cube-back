package models

import (
	"Cube-back/models/user"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func ModelRegister() {
	orm.RegisterModel(new(user.User))
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
