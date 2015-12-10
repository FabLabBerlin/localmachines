package purchases

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/models/purchases"
	"github.com/kr15h/fabsmith/tests/assert"
	"github.com/kr15h/fabsmith/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
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
	err = m.Update()
	return
}

func TestActivations(t *testing.T) {
	Convey("Testing Activation model", t, func() {

		Reset(setup.ResetDB)

		Convey("Testing CreateActivation", func() {
			user := models.User{
				FirstName: "ILoveFabLabs",
				Email:     "awesome@example.com",
			}

			Convey("Creating activation with non-existing machine", func() {
				_, err := purchases.CreateActivation(0, 0, time.Now())
				So(err, ShouldNotBeNil)
			})

			Convey("Creating activation with non-existing user", func() {
				machine, _ := CreateMachine("lel")
				_, err := purchases.CreateActivation(machine.Id, 0, time.Now())
				So(err, ShouldBeNil)
			})

			Convey("Creating activation with existing user and machine", func() {
				machine, err1 := CreateMachine("lel")
				uid, err2 := models.CreateUser(&user)
				activationStartTime := time.Date(2015, 5, 8, 2, 15, 3, 1, time.Local)
				aid, err3 := purchases.CreateActivation(machine.Id, uid, activationStartTime)
				activation, err4 := purchases.GetActivation(aid)
				assert.NoErrors(err1, err2, err3, err4)

				Convey("Activation ID should be valid", func() {
					So(aid, ShouldBeGreaterThan, 0)
				})

				Convey("It should be possible to read the activation back", func() {
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
					err := purchases.CloseActivation(aid, activationEndTime)
					assert.NoErrors(err)

					activation, err = purchases.GetActivation(aid)
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
				err := purchases.CloseActivation(0, time.Now())

				So(err, ShouldNotBeNil)
			})
			Convey("Creating an activation and close it", func() {
				machine, _ := CreateMachine("lel")
				aid, err1 := purchases.CreateActivation(machine.Id, 0, time.Now())
				err2 := purchases.CloseActivation(aid, time.Now())

				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
			})
		})
		Convey("Testing GetActivationMachineId", func() {
			Convey("Getting activation id on non-existing machine", func() {
				_, err := purchases.GetActivationMachineId(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Getting activation id on non-activated machine", func() {
				machine, err := CreateMachine("lel")
				aid, err2 := purchases.CreateActivation(machine.Id, 0, time.Now())
				err3 := purchases.CloseActivation(aid, time.Now())
				mid, err4 := purchases.GetActivationMachineId(aid)

				assert.NoErrors(err, err2, err3, err4)
				So(machine.Id, ShouldEqual, mid)
			})
			Convey("Creating activation on a machine and get activation's id", func() {
				machine, err := CreateMachine("lel")
				aid, err2 := purchases.CreateActivation(machine.Id, 0, time.Now())
				gmid, err3 := purchases.GetActivationMachineId(aid)

				assert.NoErrors(err, err2, err3)
				So(machine.Id, ShouldEqual, gmid)
			})
		})
	})
}
