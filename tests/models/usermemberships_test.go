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

func TestUserMemberships(t *testing.T) {
	Convey("Testing UserMembership model", t, func() {

		Reset(ResetDB)

		Convey("Testing CreateUserMembership", func() {

			// Create machines for the activations
			machineIdOne, _ := models.CreateMachine("Machine One")
			machineIdTwo, _ := models.CreateMachine("Machine Two")

			baseMembership := &models.Membership{}
			baseMembership.Title = "Test Membership"

			baseMembershipId, _ := models.CreateMembership(baseMembership.Title)
			baseMembership.Id = baseMembershipId
			baseMembership.DurationMonths = 1
			baseMembership.MachinePriceDeduction = 50
			baseMembership.AutoExtend = true
			baseMembership.AutoExtendDurationMonths = 30
			baseMembership.AffectedMachines = fmt.Sprintf("[%v,%v]", machineIdOne, machineIdTwo)

			models.UpdateMembership(baseMembership)
			baseMembership, _ = models.GetMembership(baseMembershipId)

			// Create a user
			user := models.User{}
			user.FirstName = "Amen"
			user.LastName = "Hesus"
			user.Email = "amen@example.com"
			userId, _ := models.CreateUser(&user)

			// Create user permissions for the created machines
			models.CreateUserPermission(userId, machineIdOne)
			models.CreateUserPermission(userId, machineIdTwo)

			// Create some activations
			timeNow := time.Date(2015, 6, 4, 0, 0, 0, 0, time.UTC)  // In membership
			timeThen := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC) // Out of membership
			CreateMembershipsActivation(userId, machineIdOne, timeNow, 5.4)
			CreateMembershipsActivation(userId, machineIdTwo, timeNow, 6.2)
			CreateMembershipsActivation(userId, machineIdOne, timeThen, 54.5)
			CreateMembershipsActivation(userId, machineIdTwo, timeThen, 12.2)

			Convey("Try creating a user membership with non existend membership ID", func() {
				fakeMembershipId := int64(-23)
				startDate := time.Now()
				var userMembershipId int64
				var err error
				fakeUserId := int64(1)
				userMembershipId, err = models.CreateUserMembership(
					fakeUserId, fakeMembershipId, startDate)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("The returned user membership ID should be 0", func() {
					So(userMembershipId, ShouldEqual, 0)
				})
			})

			//os.Exit(111)

			Convey("When creating user membership normally", func() {
				startDate := time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC)
				var err error
				var userMembershipId int64
				userMembershipId, err = models.CreateUserMembership(
					userId, baseMembershipId, startDate)
				if err != nil {
					panic(err.Error())
				}
				var gotUserMembership *models.UserMembership
				gotUserMembership, err = models.GetUserMembership(userMembershipId)
				if err != nil {
					panic(err.Error())
				}
				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("The user membership ID should be returned", func() {
					So(userMembershipId, ShouldBeGreaterThan, 0)
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
						0, int(baseMembership.DurationMonths), 0)
					So(gotUserMembership.EndDate.Equal(validEndDate), ShouldBeTrue)
				})

				Convey("The auto_extend flag should be set to the one in the base membership", func() {
					So(gotUserMembership.AutoExtend, ShouldEqual, baseMembership.AutoExtend)
				})

				Convey("The activations made during the user membership period should be affected by the base membership discount rules", func() {

					var invoice models.Invoice
					invoiceStartTime := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
					invoiceEndTime := time.Date(2015, 12, 30, 0, 0, 0, 0, time.UTC)
					invoice, err = models.CalculateInvoiceSummary(
						invoiceStartTime, invoiceEndTime)
					if err != nil {
						panic(err.Error())
					}

					// there should be 4 activations and 2 of them should be affected
					numUserSummaries := len(invoice.UserSummaries)
					So(numUserSummaries, ShouldEqual, 1)

					numActivations := len(invoice.UserSummaries[0].Activations)
					So(numActivations, ShouldEqual, 4)

					// 2 of the activations should contain memberships
					numAffectedActivations := 0
					for i := 0; i < numActivations; i++ {

						activation := invoice.UserSummaries[0].Activations[i]
						memberships := activation.Memberships
						if len(memberships) > 0 {
							numAffectedActivations += 1
						}
					}
					So(numAffectedActivations, ShouldEqual, 2)
				})
			})
		})

		Convey("Testing GetUserMembership", func() {
			Convey("Try getting a nonexistent user membership", func() {
				_, err := models.GetUserMembership(-6)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("Testing DeleteUserMembership", func() {
			baseMembership := models.Membership{}
			baseMembership.Title = "Test Membership"
			baseMembershipId, _ := models.CreateMembership(baseMembership.Title)
			baseMembership.Id = baseMembershipId

			Convey("When deleting non-existent user membership", func() {
				err := models.DeleteUserMembership(-5)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("When deleting user membership normally", func() {
				fakeUserId := int64(1)
				startDate := time.Now().UTC()
				var err error
				var userMembershipId int64
				userMembershipId, err = models.CreateUserMembership(
					fakeUserId, baseMembershipId, startDate)
				err = models.DeleteUserMembership(userMembershipId)

				Convey("It should return no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When automatically extending user membership", func() {

			// Create empty base membership
			baseMembershipId, err := models.CreateMembership("Test Membership")
			So(baseMembershipId, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			// Get the newly created base membership
			baseMembership, _ := models.GetMembership(baseMembershipId)
			//baseMembership.Duration

			// Create user membership with a start and end date some time in the past
			var userMembershipId int64
			fakeUserId := int64(1)
			loc, _ := time.LoadLocation("Europe/Berlin")
			startTime := time.Date(2015, time.July, 10, 23, 0, 0, 0, loc)

			userMembershipId, err = models.CreateUserMembership(
				fakeUserId, baseMembershipId, startTime)

			So(userMembershipId, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			// Get the created membership for later comparison
			var userMembership *models.UserMembership
			userMembership, err = models.GetUserMembership(userMembershipId)
			So(err, ShouldBeNil)
			So(userMembership, ShouldNotBeNil)

			So(userMembership.StartDate, ShouldHappenWithin,
				time.Duration(1)*time.Second, startTime)

			// Call user membership auto extend function and check the new end date
			err = models.AutoExtendUserMemberships()
			So(err, ShouldBeNil)

			//os.Exit(111)

			Convey("Check if it is extended by duration specified in the base membership", func() {

				// Get the now extended user membership
				extendedUserMembership, _ := models.GetUserMembership(userMembershipId)

				validEndDate := userMembership.EndDate.AddDate(
					0, int(baseMembership.AutoExtendDurationMonths), 0)

				So(extendedUserMembership.EndDate, ShouldHappenWithin,
					time.Duration(1)*time.Second, validEndDate)

			})
		})

	})
}
