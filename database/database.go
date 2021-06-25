package database

import (
	"Cube-back/log"
	"Cube-back/models/common/configure"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Conf struct {
	DatabaseIp       string
	DatabasePort     string
	DatabasePassword string
}

func DBValues(cmd string, args ...interface{}) (int64, []orm.Params, bool) {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(cmd, args).Values(&maps)
	if err != nil {
		log.Error(err)
		return -1, maps, false
	}
	return num, maps, true
}

func Update(s interface{}, keywords ...string) (int64, error) {
	o := orm.NewOrm()
	result, err := o.Update(s, keywords...)
	if err != nil {
		log.Error(err)
	}
	return result, err
}

func Insert(s interface{}) (int64, error) {
	o := orm.NewOrm()
	result, err := o.Insert(s)
	if err != nil {
		log.Error(err)
	}
	return result, err
}

func init() {
	conf := new(Conf)
	configure.Get(&conf)
	dataSource := "root:" + conf.DatabasePassword + "@tcp(" + conf.DatabaseIp + ":" + conf.DatabasePort + ")/cube?charset=utf8mb4"
	errDriver := orm.RegisterDriver("mysql", orm.DRMySQL)
	errDatabase := orm.RegisterDataBase("default", "mysql", dataSource)
	if errDriver != nil {
		log.Error(errDriver)
	}

	if errDatabase != nil {
		log.Error(errDatabase)
	}
	log.Info("database init successfully")
}
