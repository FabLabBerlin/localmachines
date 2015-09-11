package modelTest

import (
	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func init() {
	ConfigDB()
}

// TODO: The way go convey tests are supposed to be written is more
// human readable than this. Improve on that.

func TestMemberships(t *testing.T) {
	Convey("Testing Membership model", t, func() {
		Reset(ResetDB)

		// Base Memberships

		Convey("Testing CreateMembership", func() {
			membershipName := "Membership X"

			Convey("When creating a single membership", func() {
				membershipId, err := models.CreateMembership(membershipName)

				Convey("There should be no errors and the ID should be valid", func() {
					So(err, ShouldBeNil)
					So(membershipId, ShouldBeGreaterThan, 0)
				})

				Convey("When reading it back by using the ID", func() {
					membership, err := models.GetMembership(membershipId)

					Convey("It should return no error", func() {
						So(err, ShouldBeNil)
					})

					Convey("It should return the membership", func() {
						So(membership, ShouldNotBeNil)

						var membershipType *models.Membership
						So(membership, ShouldHaveSameTypeAs, membershipType)
					})

					Convey("Title should equal the initially given one", func() {
						So(membership.Title, ShouldEqual, membershipName)
					})

					Convey("AutoExtend should be set to true by default", func() {
						So(membership.AutoExtend, ShouldBeTrue)
					})

					Convey("AutoExtendDuration should be set to 30 by default", func() {
						So(membership.AutoExtendDuration, ShouldEqual, 30)
					})

					Convey("Unit should be set to `days`", func() {
						So(membership.Unit, ShouldEqual, "days")
					})
				})
			})
		})

		Convey("Testing GetAllMemberships", func() {

			membershipName := "The Membership"

			Convey("Getting all memberships with empty database", func() {
				memberships, err := models.GetAllMemberships()

				Convey("Should return no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should return an empty array", func() {
					So(len(memberships), ShouldEqual, 0)
				})
			})

			Convey("Getting existing memberships", func() {
				models.CreateMembership(membershipName)
				models.CreateMembership(membershipName)
				memberships, err := models.GetAllMemberships()

				Convey("Shoud return no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should return an array with 2 memberships", func() {
					So(len(memberships), ShouldEqual, 2)
				})
			})
		})

		Convey("Testing GetMembership", func() {
			membershipName := "The Membership"

			Convey("Getting non-existing membership", func() {
				_, err := models.GetMembership(0)

				Convey("Should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Creating a membership and getting it", func() {
				mid, _ := models.CreateMembership(membershipName)
				membership, err := models.GetMembership(mid)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Got membership ID should match the ID of the previously created", func() {
					So(membership.Id, ShouldEqual, mid)
				})
			})
		})

		Convey("Testing UpdateMembership", func() {
			membershipName := "Update Membership"
			newMembershipName := "New Membership Name"

			Convey("Try updating with nil object", func() {
				panicFunc := func() {
					models.UpdateMembership(nil)
				}

				Convey("There should be panic", func() {
					So(panicFunc, ShouldPanic)
				})
			})

			Convey("Try to update non existing membership", func() {
				m := &models.Membership{
					Title: membershipName,
				}
				err := models.UpdateMembership(m)

				Convey("There should be error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Create membership and update it", func() {
				mid, _ := models.CreateMembership(membershipName)
				m, _ := models.GetMembership(mid)
				m.Title = newMembershipName
				err := models.UpdateMembership(m)
				nm, _ := models.GetMembership(mid)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("The title of the read-back membership should equal given one", func() {
					So(nm.Title, ShouldEqual, newMembershipName)
				})
			})
		})

		Convey("Testing DeleteMembership", func() {
			membershipName := "Super Membership"

			Convey("Try to delete non-existing membership", func() {
				err := models.DeleteMembership(0)

				Convey("It should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Creating a membership and delete it", func() {
				mid, _ := models.CreateMembership(membershipName)
				err := models.DeleteMembership(mid)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		// User Memberships

		Convey("Testing CreateUserMembership", func() {
			baseMembership := &models.Membership{}
			baseMembership.Title = "Test Membership"
			baseMembershipId, _ := models.CreateMembership(baseMembership.Title)
			baseMembership, _ = models.GetMembership(baseMembershipId)

			Convey("Try creating a user membership with non existend membership ID", func() {
				fakeUserId := int64(1)
				fakeMembershipId := int64(-23)
				startDate := time.Now()
				var userMembershipId int64
				var err error
				userMembershipId, err = models.CreateUserMembership(
					fakeUserId, fakeMembershipId, startDate)

				Convey("It should return error", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("The returned user membership ID should be 0", func() {
					So(userMembershipId, ShouldEqual, 0)
				})
			})

			Convey("When creating user membership normally", func() {
				fakeUserId := int64(1)
				startDate := time.Now().UTC()
				var err error
				var userMembershipId int64
				userMembershipId, err = models.CreateUserMembership(
					fakeUserId, baseMembershipId, startDate)
				var gotUserMembership *models.UserMembership
				gotUserMembership, err = models.GetUserMembership(userMembershipId)

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

				Convey("The end date should be correct according to the base membership", func() {
					validEndDate := gotUserMembership.StartDate.AddDate(0, 0, int(baseMembership.Duration))
					So(gotUserMembership.EndDate.Equal(validEndDate), ShouldBeTrue)
				})

				Convey("The auto_extend flag should be set to the one in the base membership", func() {
					So(gotUserMembership.AutoExtend, ShouldEqual, baseMembership.AutoExtend)
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

			// Create user membership with a start and end date some time in the past
			var userMembershipId int64
			fakeUserId := int64(1)
			startDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
			userMembershipId, err = models.CreateUserMembership(
				fakeUserId, baseMembershipId, startDate)
			So(userMembershipId, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			// Get the created membership for later comparison
			userMembership, _ := models.GetUserMembership(userMembershipId)
			So(userMembership, ShouldNotBeNil)

			// Call user membership auto extend function and check the new end date
			err = models.AutoExtendUserMemberships()
			So(err, ShouldBeNil)

			Convey("Check if it is extended by duration specified in the base membership", func() {

				// Get the now extended user membership
				extendedUserMembership, _ := models.GetUserMembership(userMembershipId)
				validEndDate := userMembership.EndDate.AddDate(0, 0, int(baseMembership.AutoExtendDuration))
				So(extendedUserMembership.EndDate.Equal(validEndDate), ShouldBeTrue)

			})
		})
	})
}
