package userTests

import (
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/kr15h/fabsmith/tests/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestUsers(t *testing.T) {
	SkipConvey("Testing Activation model", t, func() {
		Reset(func() {
			o := orm.NewOrm()
			var activations []models.Activation
			o.QueryTable("activations").All(&activations)
			for _, item := range activations {
				o.Delete(&item)
			}
		})
	})
}
