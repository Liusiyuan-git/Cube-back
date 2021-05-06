package database

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int
	Email    string
	Password string
	Phone    string
}

func DBValues(cmd string, args ...interface{}) (int64, []orm.Params, bool) {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(cmd, args).Values(&maps)
	if err != nil {
		fmt.Println(err)
		return -1, maps, false
	}
	return num, maps, true
}

func Update(s interface{}, keywords string) (int64, error) {
	o := orm.NewOrm()
	result, err := o.Update(s, keywords)
	if err != nil {
		logs.Error(err)
	}
	return result, err
}

func Insert(s interface{}) (int64, error) {
	o := orm.NewOrm()
	result, err := o.Insert(s)
	if err != nil {
		logs.Error(err)
	}
	return result, err
}

func init() {
	errDriver := orm.RegisterDriver("mysql", orm.DRMySQL)
	errDatabase := orm.RegisterDataBase("default", "mysql", "root:201020120402ssS~@tcp(81.68.121.120:3306)/cube?charset=utf8")
	if errDriver != nil {
		logs.Error(errDriver)
	}

	if errDatabase != nil {
		logs.Error(errDatabase)
	}
}
