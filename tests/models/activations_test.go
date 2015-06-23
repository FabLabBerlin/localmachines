package modelTest

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestActivations(t *testing.T) {
	Convey("Testing Activation model", t, func() {
		Reset(ResetDB)
		Convey("Testing CreateActivation", func() {
			user := models.User{FirstName: "ILoveFabLabs"}
			Convey("Creating activation with non-existing machine", func() {
				_, err := models.CreateActivation(0, 0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating activation with non-existing user", func() {
				mid, _ := models.CreateMachine("lel")
				_, err := models.CreateActivation(mid, 0)

				So(err, ShouldBeNil)
			})
			Convey("Creating activation with existing user and machine", func() {
				mid, _ := models.CreateMachine("lel")
				uid, _ := models.CreateUser(&user)
				_, err := models.CreateActivation(mid, uid)

				So(err, ShouldBeNil)
			})
		})
	})
}
