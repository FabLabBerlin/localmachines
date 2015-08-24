package controllerTest

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/routers"
	"github.com/kr15h/fabsmith/tests/models"
	"os"
)

func init() {
	modelTest.ConfigDB()
	controllerTestsDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	beego.TestBeegoInit(controllerTestsDir)
	routers.Init()
}
