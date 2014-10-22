package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
)

func init() {
	// Template init
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"

	// Database init
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/fabsmith?charset=utf8")
}

func main() {
	beego.Run()
}
