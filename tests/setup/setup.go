package setup

import (
	"github.com/astaxie/beego"
	"os"
)

func ConfigDir(dirSuffix string) {
	modelTestsDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	beego.TestBeegoInit(modelTestsDir + dirSuffix)
}
