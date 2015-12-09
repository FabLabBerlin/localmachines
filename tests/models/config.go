package modelTest

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// Beego ORM model used only for testing time
type TimeTest struct {
	Id   int64     `orm:"auto";"pk"`
	Time time.Time `orm:"type(timestamp)"`
}

func (this *TimeTest) TableName() string {
	return "time_test"
}

func init() {
	orm.RegisterModel(new(TimeTest))
}
