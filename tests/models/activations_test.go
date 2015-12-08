package modelTest

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func CreateMachine(name string) (m *models.Machine, err error) {
	mid, err := models.CreateMachine(name)
	if err != nil {
		return
	}
	if m, err = models.GetMachine(mid); err != nil {
		return
	}
	m.Price = 0.1
	m.PriceUnit = "minute"
	err = models.UpdateMachine(m)
	return
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
				machine, _ := CreateMachine("lel")
				_, err := models.CreateActivation(machine.Id, 0, time.Now())
				So(err, ShouldBeNil)
			})

			Convey("Creating activation with existing user and machine", func() {
				machine, _ := CreateMachine("lel")
				uid, _ := models.CreateUser(&user)
				activationStartTime := time.Date(2015, 5, 8, 2, 15, 3, 1, time.Local)
				aid, err := models.CreateActivation(machine.Id, uid, activationStartTime)
				if err != nil {
					panic(fmt.Sprintf("create activation: %v", err))
				}
				activation, err2 := models.GetActivation(aid)
				if err2 != nil {
					panic(fmt.Sprintf("get activation: %v", err2))
				}

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Activation ID should be valid", func() {
					So(aid, ShouldBeGreaterThan, 0)
				})

				Convey("It should be possible to read the activation back", func() {
					So(err2, ShouldBeNil)
					So(activation.Purchase.Id, ShouldEqual, aid)
				})

				Convey("Activation start time should match current time", func() {
					So(activation.Purchase.TimeStart, ShouldHappenWithin,
						time.Duration(1)*time.Second, activationStartTime)
				})

				Convey("the active flag should be true after creating", func() {
					So(activation.Purchase.Running, ShouldBeTrue)
				})

				Convey("It should be possible to close the activation", func() {
					activationEndTime := time.Now()
					err = models.CloseActivation(aid, activationEndTime)
					if err != nil {
						panic(fmt.Sprintf("close activation: %v", err))
					}
					So(err, ShouldBeNil)

					activation, err = models.GetActivation(aid)
					if err != nil {
						panic(fmt.Sprintf("get activation: %v", err))
					}

					Convey("The end time of the closed activation should be correct", func() {
						So(activation.Purchase.TimeEnd, ShouldHappenWithin,
							time.Duration(1)*time.Second, activationEndTime)
					})

					Convey("The total duration of the activation should be correct", func() {
						totalTime := activation.Purchase.TimeEnd.Sub(activation.Purchase.TimeStart)
						q := activation.Purchase.Quantity
						So(q, ShouldAlmostEqual,
							int64(totalTime.Minutes()), 1)
					})

					Convey("the active flag should be false after closing", func() {
						So(activation.Purchase.Running, ShouldBeFalse)
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
				machine, _ := CreateMachine("lel")
				aid, _ := models.CreateActivation(machine.Id, 0, time.Now())
				err := models.CloseActivation(aid, time.Now())

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing GetActivationMachineId", func() {
			Convey("Getting activation id on non-existing machine", func() {
				_, err := models.GetActivationMachineId(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Getting activation id on non-activated machine", func() {
				machine, _ := CreateMachine("lel")
				aid, _ := models.CreateActivation(machine.Id, 0, time.Now())
				models.CloseActivation(aid, time.Now())
				_, err := models.GetActivationMachineId(aid)

				So(err, ShouldBeNil)
			})
			Convey("Creating activation on a machine and get activation's id", func() {
				machine, _ := CreateMachine("lel")
				aid, _ := models.CreateActivation(machine.Id, 0, time.Now())
				gmid, _ := models.GetActivationMachineId(aid)

				So(machine.Id, ShouldEqual, gmid)
			})
		})
	})
}
