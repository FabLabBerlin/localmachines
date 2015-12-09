package controllerTest

import (
	"github.com/kr15h/fabsmith/routers"
	"github.com/kr15h/fabsmith/tests/setup"
)

func init() {
	setup.ConfigDir("/..")
	setup.ConfigDB()
	routers.Init()
}
