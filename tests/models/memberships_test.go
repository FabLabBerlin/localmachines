package modelTest

import (
	"testing"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestMemberships(t *testing.T) {
	Convey("Testing Membership model", t, func() {
		Reset(ResetDB)
		Convey("Testing GetAllMemberships", func() {
			membershipName := "Lel"
			Convey("Getting all memberships when there is nothing in the database", func() {
				memberships, err := models.GetAllMemberships()

				So(err, ShouldBeNil)
				So(len(memberships), ShouldEqual, 0)
			})
			Convey("Creating 2 memberships and get them all", func() {
				models.CreateMembership(membershipName)
				models.CreateMembership(membershipName)
				memberships, err := models.GetAllMemberships()

				So(err, ShouldBeNil)
				So(len(memberships), ShouldEqual, 2)
			})
		})
		Convey("Testing CreateMembership", func() {
			membershipName := "My awesome membership"
			Convey("Creating one membership", func() {
				_, err := models.CreateMembership(membershipName)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing GetMembership", func() {
			membershipName := "Lel"
			Convey("Getting non-existing membership", func() {
				_, err := models.GetMembership(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating a membership and getting it", func() {
				mid, _ := models.CreateMembership(membershipName)
				membership, err := models.GetMembership(mid)

				So(membership.Id, ShouldEqual, mid)
				So(err, ShouldBeNil)
			})
		})
		Convey("Testing UpdateMembership", func() {
			membershipName := "Lel"
			newMembershipName := "DatAwesomeNewName"
			Convey("Try updating with nil object", func() {
				panicFunc := func() {
					models.UpdateMembership(nil)
				}

				So(panicFunc, ShouldPanic)
			})
			Convey("Try to update non existing membership", func() {
				m := &models.Membership{
					Title: membershipName,
				}
				err := models.UpdateMembership(m)

				So(err, ShouldNotBeNil)
			})
			Convey("Create membership and update it", func() {
				mid, _ := models.CreateMembership(membershipName)
				m, _ := models.GetMembership(mid)
				m.Title = newMembershipName
				err := models.UpdateMembership(m)
				nm, _ := models.GetMembership(mid)

				So(err, ShouldBeNil)
				So(nm.Title, ShouldEqual, newMembershipName)
			})
		})
		Convey("Testing DeleteMembership", func() {
			membershipName := "My awesome membership program"
			Convey("Try to delete non-existing membership", func() {
				err := models.DeleteMembership(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating a membership and delete it", func() {
				mid, _ := models.CreateMembership(membershipName)
				err := models.DeleteMembership(mid)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing CreateUserMembership", func() {
			Convey("Try creating a user with nil parameter", func() {
				_, err := models.CreateUserMembership(nil)

				So(err, ShouldNotBeNil)
			})
		})
	})
}
