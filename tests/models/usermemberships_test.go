package modelTest

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/memberships/auto_extend"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
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

			membership, err := memberships.CreateMembership(1, "Test Membership")
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
			if membership, err = memberships.GetMembership(membership.Id); err != nil {
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

			// Create some activations
			timeNow := time.Date(2015, 6, 4, 0, 0, 0, 0, time.UTC)  // In membership
			timeThen := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC) // Out of membership
			CreateMembershipsActivation(userId, machineOne.Id, timeNow, 5.4)
			CreateMembershipsActivation(userId, machineTwo.Id, timeNow, 6.2)
			CreateMembershipsActivation(userId, machineOne.Id, timeThen, 54.5)
			CreateMembershipsActivation(userId, machineTwo.Id, timeThen, 12.2)

			Convey("Try creating a user membership with non existend membership ID", func() {
				fakeMembershipId := int64(-23)
				startDate := time.Now()
				fakeUserId := int64(1)
				userMembership, err := memberships.CreateUserMembership(
					fakeUserId, fakeMembershipId, startDate)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("The returned user membership should be nil", func() {
					So(userMembership, ShouldBeNil)
				})
			})

			Convey("When creating user membership normally", func() {
				startDate := time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC)
				userMembership, err := memberships.CreateUserMembership(
					userId, membership.Id, startDate)
				if err != nil {
					panic(err.Error())
				}
				gotUserMembership, err := memberships.GetUserMembership(userMembership.Id)
				if err != nil {
					panic(err.Error())
				}
				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("The user membership ID should be returned", func() {
					So(userMembership.Id, ShouldBeGreaterThan, 0)
				})

				Convey("It should be possible to read it back again", func() {
					So(err, ShouldBeNil)
					So(gotUserMembership, ShouldNotBeNil)
				})

				Convey("The start date should be correct", func() {
					So(gotUserMembership.StartDate, ShouldHappenWithin,
						time.Duration(1)*time.Second, startDate)
				})

				Convey("The end date should be correct according to the base membership", func() {
					validEndDate := gotUserMembership.StartDate.AddDate(
						0, int(membership.DurationMonths), 0)
					So(gotUserMembership.EndDate.Equal(validEndDate), ShouldBeTrue)
				})

				Convey("The auto_extend flag should be set to the one in the base membership", func() {
					So(gotUserMembership.AutoExtend, ShouldEqual, membership.AutoExtend)
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
				_, err := memberships.GetUserMembership(-6)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When automatically extending user membership", func() {
			// Create empty base membership
			m, err := memberships.CreateMembership(1, "Test Membership")
			So(m.Id, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			// Create user membership with a start and end date some time in the past
			fakeUserId := int64(1)
			loc, _ := time.LoadLocation("Europe/Berlin")
			startTime := time.Date(2015, time.July, 10, 23, 0, 0, 0, loc)

			userMembership, err := memberships.CreateUserMembership(
				fakeUserId, m.Id, startTime)

			So(userMembership.Id, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			inv := &invoices.Invoice{
				LocationId: 1,
				UserId:     fakeUserId,
				Month:      int(time.Now().Month()),
				Year:       time.Now().Year(),
			}
			if _, err = invoices.CreateOrUpdate(inv); err != nil {
				panic(err.Error())
			}

			// Get the created membership for later comparison
			userMembership, err = memberships.GetUserMembership(userMembership.Id)
			So(err, ShouldBeNil)
			So(userMembership, ShouldNotBeNil)

			So(userMembership.StartDate, ShouldHappenWithin,
				time.Duration(1)*time.Second, startTime)

			if err = inv.AttachUserMembership(userMembership); err != nil {
				panic(err.Error())
			}

			// Call user membership auto extend function and check the new end date
			err = auto_extend.AutoExtendUserMemberships()
			if err != nil {
				panic(err.Error())
			}

			Convey("Check if it is extended by duration specified in the base membership", func() {

				// Get the now extended user membership
				extendedUserMembership, _ := memberships.GetUserMembership(userMembership.Id)

				validEndDate := userMembership.EndDate.AddDate(
					0, int(m.AutoExtendDurationMonths), 0)

				So(extendedUserMembership.EndDate, ShouldHappenWithin,
					time.Duration(1)*time.Second, validEndDate)

			})
		})

		Convey("When automatically extending an already extended user membership", func() {
			Reset(setup.ResetDB)

			m, err := memberships.CreateMembership(1, "Test Membership")
			if err != nil {
				panic(err.Error())
			}
			So(m.AutoExtend, ShouldBeTrue)
			So(m.DurationMonths, ShouldEqual, 1)
			So(m.AutoExtendDurationMonths, ShouldEqual, 1)

			// Create user membership with a start and end date some time in the past
			fakeUserId := int64(11)
			loc, _ := time.LoadLocation("Europe/Berlin")
			startTime := time.Date(2015, time.July, 10, 23, 0, 0, 0, loc)

			userMembership, err := memberships.CreateUserMembership(
				fakeUserId, m.Id, startTime)
			inv7_15 := &invoices.Invoice{
				LocationId: 1,
				UserId:     fakeUserId,
				Month:      int(time.July),
				Year:       2015,
			}
			if _, err = invoices.CreateOrUpdate(inv7_15); err != nil {
				panic(err.Error())
			}

			if err = inv7_15.AttachUserMembership(userMembership); err != nil {
				panic(err.Error())
			}

			invNow := &invoices.Invoice{
				LocationId: 1,
				UserId:     fakeUserId,
				Month:      int(time.Now().Month()),
				Year:       time.Now().Year(),
			}
			if _, err = invoices.CreateOrUpdate(invNow); err != nil {
				panic(err.Error())
			}

			if err = invNow.AttachUserMembership(userMembership); err != nil {
				panic(err.Error())
			}

			// Call user membership auto extend function and check the new end date
			err = auto_extend.AutoExtendUserMemberships()
			if err != nil {
				panic(err.Error())
			}

			ums, err := memberships.GetAllUserMembershipsAt(1)
			if err != nil {
				panic(err.Error())
			}
			So(len(ums), ShouldEqual, 2)
			um7_15 := ums[0]
			umNow := ums[1]
			So(um7_15.InvoiceId, ShouldEqual, inv7_15.Id)
			So(umNow.InvoiceId, ShouldEqual, invNow.Id)

			So(um7_15.StartDate.Format("2006-01"), ShouldEqual, startTime.Format("2006-01"))
			So(um7_15.EndDate.Format("2006-01"), ShouldEqual, startTime.AddDate(0, 1, 0).Format("2006-01"))

			So(umNow.StartDate.Format("2006-01"), ShouldEqual, startTime.Format("2006-01"))
			So(umNow.EndDate.Format("2006-01"), ShouldEqual, startTime.AddDate(0, 2, 0).Format("2006-01"))

			// Call auto extend user membership often enough, so it extends
			// until today
			for i := 0; i < 100; i++ {
				err = auto_extend.AutoExtendUserMemberships()
				if err != nil {
					panic(err.Error())
				}
			}

			umNow, err = memberships.GetUserMembership(umNow.Id)
			if err != nil {
				panic(err.Error())
			}
			So(umNow.StartDate.Format("2006-01"), ShouldEqual, startTime.Format("2006-01"))
			So(umNow.EndDate.Format("2006-01"), ShouldEqual, time.Now().Format("2006-01"))

			ums, err = memberships.GetAllUserMembershipsAt(1)
			if err != nil {
				panic(err.Error())
			}
			So(len(ums), ShouldEqual, 2)
		})

	})
}
