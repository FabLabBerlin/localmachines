package connect

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"net/url"
)

func Connect(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb string) {
	loc := url.QueryEscape("UTC")
	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb, loc)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}
