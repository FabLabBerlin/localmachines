package modelTest

import (
	"testing"

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
		Convey("Testing CloseActivation", func() {
			Convey("Trying to close a non-existing activation", func() {
				err := models.CloseActivation(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating an activation and close it", func() {
				mid, _ := models.CreateMachine("lel")
				aid, _ := models.CreateActivation(mid, 0)
				err := models.CloseActivation(aid)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing DeleteActivation", func() {
			Convey("Trying to delete non-existing activation", func() {
				err := models.DeleteActivation(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating an activation and delete it", func() {
				mid, _ := models.CreateMachine("lel")
				aid, _ := models.CreateActivation(mid, 0)
				err := models.DeleteActivation(aid)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing GetActivationMachineId", func() {
			Convey("Getting activation id on non-existing machine", func() {
				_, err := models.GetActivationMachineId(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Getting activation id on non-activated machine", func() {
				mid, _ := models.CreateMachine("lel")
				aid, _ := models.CreateActivation(mid, 0)
				models.CloseActivation(aid)
				_, err := models.GetActivationMachineId(aid)

				So(err, ShouldBeNil)
			})
			Convey("Creating activation on a machine and get activation's id", func() {
				mid, _ := models.CreateMachine("lel")
				aid, _ := models.CreateActivation(mid, 0)
				gmid, _ := models.GetActivationMachineId(aid)

				So(mid, ShouldEqual, gmid)
			})
		})
	})
}
