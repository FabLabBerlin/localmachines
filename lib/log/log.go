package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func init() {
	logs.Register("logrus", func() logs.Logger {
		return Logger{}
	})
	if err := beego.SetLogger("logrus", ""); err != nil {
		panic(err.Error())
	}
}

type Logger struct {
	logs.Logger
}

// Adapter to beego's Logger interface
func (l Logger) Init(config string) error {
	return nil
}

func (l Logger) WriteMsg(msg string, level int) error {
	switch level {
	case 0:
		logrus.Panicln(msg)
	case 1:
		logrus.Fatalln(msg)
	case 2:
		logrus.Errorln(msg)
	case 3:
		logrus.Warningln(msg)
	case 4:
		logrus.Infoln(msg)
	default:
		logrus.Debugln(msg)
	}
	return nil
}

func (l Logger) Destroy() {
	return
}

func (l Logger) Flush() {
	return
}
