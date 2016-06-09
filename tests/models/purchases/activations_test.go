package purchases

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/assert"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func CreateMachine(name string) (m *machine.Machine, err error) {
	if m, err = machine.Create(1, name); err != nil {
		return
	}
	m.Price = 0.1
	m.PriceUnit = "minute"
	err = m.Update(false)
	return
}

func TestActivations(t *testing.T) {
	Convey("Testing Activation model", t, func() {

		Reset(setup.ResetDB)

		Convey("Testing StartActivation", func() {
			user := users.User{
				FirstName: "ILoveFabLabs",
				Email:     "awesome@example.com",
			}

			Convey("Creating activation with non-existing user", func() {
				machine, _ := CreateMachine("lel")
				_, err := purchases.StartActivation(machine, 0, time.Now())
				So(err, ShouldBeNil)
			})

			Convey("Creating activation with existing user and machine", func() {
				machine, err1 := CreateMachine("lel")
				uid, err2 := users.CreateUser(&user)
				activationStartTime := time.Date(2015, 5, 8, 2, 15, 3, 1, time.Local)
				aid, err3 := purchases.StartActivation(machine, uid, activationStartTime)
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
					a, err := purchases.GetActivation(aid)
					if err != nil {
						panic(err.Error())
					}
					err = a.Close(activationEndTime)
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

					Convey("Starting is idempotent", func() {
						machine, err := CreateMachine("lel")
						if err != nil {
							panic(err.Error())
						}
						uid, err := users.CreateUser(&users.User{
							FirstName: "Ron Sommer",
							Email:     "ron@dtag.de",
						})
						if err != nil {
							panic(err.Error())
						}
						activationStartTime := time.Date(2015, 5, 8, 2, 15, 3, 1, time.Local)
						aid, err := purchases.StartActivation(machine, uid, activationStartTime)
						if err != nil {
							panic(err.Error())
						}
						_ = aid
						aid2, err := purchases.StartActivation(machine, uid, activationStartTime)
						if err != nil {
							panic(err.Error())
						}
						_ = aid2
						So(aid, ShouldEqual, aid2)
					})
				})
			})
		})

		Convey("Testing CloseActivation", func() {
			Convey("Trying to close a non-existing activation", func() {
				_, err := purchases.GetActivation(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Starting an activation and closing it", func() {
				machine, _ := CreateMachine("lel")
				aid, err1 := purchases.StartActivation(machine, 0, time.Now())
				a, err := purchases.GetActivation(aid)
				if err != nil {
					panic(err.Error())
				}
				err2 := a.Close(time.Now())

				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
			})

			Convey("Closing is idempotent", func() {
				machine, err := CreateMachine("lel")
				if err != nil {
					panic(err.Error())
				}
				aid, err := purchases.StartActivation(machine, 0, time.Now())
				if err != nil {
					panic(err.Error())
				}
				a, err := purchases.GetActivation(aid)
				if err != nil {
					panic(err.Error())
				}
				err = a.Close(time.Now())
				if err != nil {
					panic(err.Error())
				}
				err = a.Close(time.Now())
				So(err, ShouldBeNil)
			})
		})
	})
}
