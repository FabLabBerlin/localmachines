package modelTest

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/invoices/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestUserMemberships(t *testing.T) {
	Convey("Testing UserMembership model", t, func() {

		Reset(setup.ResetDB)

		Convey("Testing CreateUserMembership", func() {

			// Create machines for the activations
			machineOne, _ := machine.Create(1, "Machine One")
			machineTwo, _ := machine.Create(1, "Machine Two")

			membership, err := memberships.Create(1, "Test Membership")
			if err != nil {
				panic(err.Error())
			}
			membership.DurationMonths = 1
			membership.MachinePriceDeduction = 50
			membership.AutoExtend = true
			membership.AutoExtendDurationMonths = 30
			membership.AffectedMachines = fmt.Sprintf("[%v,%v]", machineOne.Id, machineTwo.Id)
			if err := membership.Update(); err != nil {
				panic(err.Error())
			}
			if membership, err = memberships.Get(membership.Id); err != nil {
				panic(err.Error())
			}

			// Create a user
			user := users.User{
				FirstName: "Amen",
				LastName:  "Hesus",
				Email:     "amen@example.com",
			}
			userId, _ := users.CreateUser(&user)
			_, err = user_locations.Create(&user_locations.UserLocation{
				UserId:     user.Id,
				LocationId: 1,
			})
			if err != nil {
				panic(err.Error())
			}

			// Create user permissions for the created machines
			user_permissions.Create(userId, machineOne.Id)
			user_permissions.Create(userId, machineTwo.Id)

			invNow := &invutil.Invoice{}
			invNow.LocationId = 1
			invNow.UserId = userId
			invNow.Month = 6
			invNow.Year = 2015
			invNow.Status = "draft"
			if _, err = invoices.Create(&invNow.Invoice); err != nil {
				panic(err.Error())
			}

			invThen := &invutil.Invoice{}
			invThen.LocationId = 1
			invThen.UserId = userId
			invThen.Month = 2
			invThen.Year = 2015
			invThen.Status = "draft"
			if _, err = invoices.Create(&invThen.Invoice); err != nil {
				panic(err.Error())
			}

			// Create some activations
			timeNow := time.Date(2015, 6, 4, 0, 0, 0, 0, time.UTC)  // In membership
			timeThen := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC) // Out of membership
			CreateMembershipsActivation(userId, machineOne.Id, invNow.Id, timeNow, 5.4)
			CreateMembershipsActivation(userId, machineTwo.Id, invNow.Id, timeNow, 6.2)
			CreateMembershipsActivation(userId, machineOne.Id, invThen.Id, timeThen, 54.5)
			CreateMembershipsActivation(userId, machineTwo.Id, invThen.Id, timeThen, 12.2)

			Convey("Try creating a user membership with non existend membership ID", func() {
				fakeMembershipId := int64(-23)
				startDate := time.Now()
				fakeUserId := int64(1)

				o := orm.NewOrm()

				userMembership, err := user_memberships.Create(
					o, fakeUserId, fakeMembershipId, 123, startDate)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("The returned user membership should be nil", func() {
					So(userMembership, ShouldBeNil)
				})
			})

			Convey("When creating user membership normally", func() {
				startDate := time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC)

				o := orm.NewOrm()

				userMembership, err := user_memberships.Create(
					o, userId, membership.Id, invNow.Id, startDate)
				if err != nil {
					panic(err.Error())
				}

				o = orm.NewOrm()

				if err := userMembership.Update(o); err != nil {
					panic(err.Error())
				}

				gotUserMembership, err := user_memberships.Get(userMembership.Id)
				if err != nil {
					panic(err.Error())
				}

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("The user membership ID should be returned", func() {
					So(userMembership.Id, ShouldBeGreaterThan, 0)
				})

				Convey("The combo type should work normally", func() {
					ums, err := user_memberships.GetAllAt(1)
					if err != nil {
						panic(err.Error())
					}
					So(len(ums), ShouldEqual, 1)
				})

				Convey("It should be possible to read it back again", func() {
					So(err, ShouldBeNil)
					So(gotUserMembership, ShouldNotBeNil)
				})

				Convey("The start date should be correct", func() {
					So(gotUserMembership.StartDate, ShouldHappenWithin,
						time.Duration(1)*time.Second, startDate)
				})

				Convey("The activations made during the user membership period should be affected by the base membership discount rules", func() {
					interval := lib.Interval{
						MonthFrom: 6,
						YearFrom:  2015,
						MonthTo:   6,
						YearTo:    2015,
					}

					me, err := monthly_earning.New(1, interval)
					if err != nil {
						panic(err.Error())
					}

					// there should be 2 activations and 2 of them should be affected
					numUserSummaries := len(me.Invoices)
					So(numUserSummaries, ShouldEqual, 1)

					numActivations := len(me.Invoices[0].Purchases)
					So(numActivations, ShouldEqual, 2)

					// 2 of the activations should contain memberships
					numAffectedActivations := 0
					for i := 0; i < numActivations; i++ {

						activation := me.Invoices[0].Purchases[i]
						memberships := activation.Memberships
						if len(memberships) > 0 {
							numAffectedActivations += 1
						}
					}
					So(numAffectedActivations, ShouldEqual, 2)
				})
			})
		})

		Convey("GetUserMembership", func() {
			Convey("Try getting a nonexistent user membership", func() {
				_, err := user_memberships.Get(-6)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
