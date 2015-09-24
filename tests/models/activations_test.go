package modelTest

import (
	"testing"
	"time"

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
				_, err := models.CreateActivation(0, 0, time.Now())
				So(err, ShouldNotBeNil)
			})

			Convey("Creating activation with non-existing user", func() {
				mid, _ := models.CreateMachine("lel")
				_, err := models.CreateActivation(mid, 0, time.Now())
				So(err, ShouldBeNil)
			})

			Convey("Creating activation with existing user and machine", func() {
				mid, _ := models.CreateMachine("lel")
				uid, _ := models.CreateUser(&user)
				activationStartTime := time.Date(2015, 5, 8, 2, 15, 3, 1, time.Local)
				aid, err := models.CreateActivation(mid, uid, activationStartTime)
				activation, err2 := models.GetActivation(aid)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Activation ID should be valid", func() {
					So(aid, ShouldBeGreaterThan, 0)
				})

				Convey("It should be possible to read the activation back", func() {
					So(err2, ShouldBeNil)
					So(activation.Id, ShouldEqual, aid)
				})

				Convey("Activation start time should match current time", func() {
					So(activation.TimeStart, ShouldHappenWithin,
						time.Duration(1)*time.Second, activationStartTime)
				})

				Convey("the active flag should be true after creating", func() {
					So(activation.Active, ShouldBeTrue)
				})

				Convey("It should be possible to close the activation", func() {
					activationEndTime := time.Now()
					err = models.CloseActivation(aid, activationEndTime)
					So(err, ShouldBeNil)

					activation, _ = models.GetActivation(aid)

					Convey("The end time of the closed activation should be correct", func() {
						So(activation.TimeEnd, ShouldHappenWithin,
							time.Duration(1)*time.Second, activationEndTime)
					})

					Convey("The total duration of the activation should be correct", func() {
						totalTime := activation.TimeEnd.Sub(activation.TimeStart)
						So(activation.TimeTotal, ShouldAlmostEqual,
							int64(totalTime.Seconds()), 1)
					})

					Convey("the active flag should be false after closing", func() {
						So(activation.Active, ShouldBeFalse)
					})
				})
			})
		})

		Convey("Testing CloseActivation", func() {
			Convey("Trying to close a non-existing activation", func() {
				err := models.CloseActivation(0, time.Now())

				So(err, ShouldNotBeNil)
			})
			Convey("Creating an activation and close it", func() {
				mid, _ := models.CreateMachine("lel")
				aid, _ := models.CreateActivation(mid, 0, time.Now())
				err := models.CloseActivation(aid, time.Now())

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
				aid, _ := models.CreateActivation(mid, 0, time.Now())
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
				aid, _ := models.CreateActivation(mid, 0, time.Now())
				models.CloseActivation(aid, time.Now())
				_, err := models.GetActivationMachineId(aid)

				So(err, ShouldBeNil)
			})
			Convey("Creating activation on a machine and get activation's id", func() {
				mid, _ := models.CreateMachine("lel")
				aid, _ := models.CreateActivation(mid, 0, time.Now())
				gmid, _ := models.GetActivationMachineId(aid)

				So(mid, ShouldEqual, gmid)
			})
		})
	})
}
