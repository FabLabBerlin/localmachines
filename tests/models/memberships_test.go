package modelTest

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

// TODO: The way go convey tests are supposed to be written is more
// human readable than this. Improve on that.

func CreateMembershipsActivation(userId, machineId int64, startTime time.Time, minutes float64) (id int64) {

	activation := purchases.Activation{
		Purchase: purchases.Purchase{
			LocationId: 1,
			TimeStart:  startTime,
			UserId:     userId,
			MachineId:  machineId,
		},
	}
	activation.Purchase.TimeEnd = activation.Purchase.TimeStart.Add(time.Duration(minutes) * time.Minute)

	o := orm.NewOrm()
	id, err := o.Insert(&activation.Purchase)
	if err != nil {
		panic(err.Error())
	}
	return
}

func TestMemberships(t *testing.T) {
	Convey("Testing Membership model", t, func() {

		Reset(setup.ResetDB)

		Convey("Testing CreateMembership", func() {
			membershipName := "Membership X"

			Convey("When creating a single membership", func() {
				membership, err := models.CreateMembership(1, membershipName)
				if err != nil {
					panic(err.Error())
				}

				Convey("There should be no errors and the ID should be valid", func() {
					So(err, ShouldBeNil)
					So(membership.Id, ShouldBeGreaterThan, 0)
				})

				Convey("When reading it back by using the ID", func() {
					membership, err := models.GetMembership(membership.Id)
					if err != nil {
						panic(fmt.Sprintf("%v ... membershipId: %v", err.Error(), membership.Id))
					}
					Convey("It should return no error", func() {
						So(err, ShouldBeNil)
					})

					Convey("It should return the membership", func() {
						So(membership, ShouldNotBeNil)

						var membershipType *models.Membership
						So(membership, ShouldHaveSameTypeAs, membershipType)
					})

					Convey("Title should equal the initially given one", func() {
						title := membership.Title
						So(title, ShouldEqual, membershipName)
					})

					Convey("The duration of the membership should be set "+
						"to 1 month by default", func() {

						So(membership.DurationMonths, ShouldEqual, 1)
					})

					Convey("AutoExtend should be set to true by default", func() {
						So(membership.AutoExtend, ShouldBeTrue)
					})

					Convey("AutoExtendDuration in months should be set to 1 by default", func() {
						So(membership.AutoExtendDurationMonths, ShouldEqual, 1)
					})
				})
			})
		})

		Convey("Testing GetAllMemberships", func() {

			membershipName := "The Membership"

			Convey("Getting all memberships with empty database", func() {
				memberships, err := models.GetAllMembershipsAt(1)

				Convey("Should return no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should return an empty array", func() {
					So(len(memberships), ShouldEqual, 0)
				})
			})

			Convey("Getting existing memberships", func() {
				models.CreateMembership(1, membershipName)
				models.CreateMembership(1, membershipName)
				memberships, err := models.GetAllMembershipsAt(1)

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
				m, _ := models.CreateMembership(1, membershipName)
				membership, err := models.GetMembership(m.Id)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Got membership ID should match the ID of the previously created", func() {
					So(membership.Id, ShouldEqual, m.Id)
				})
			})
		})

		Convey("Testing UpdateMembership", func() {
			membershipName := "Update Membership"
			newMembershipName := "New Membership Name"

			Convey("Try to update non existing membership", func() {
				m := &models.Membership{
					Title: membershipName,
				}
				err := m.Update()

				Convey("There should be error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Create membership and update it", func() {
				m, _ := models.CreateMembership(1, membershipName)
				m.Title = newMembershipName
				err := m.Update()
				nm, _ := models.GetMembership(m.Id)

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
				m, _ := models.CreateMembership(1, membershipName)
				err := models.DeleteMembership(m.Id)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

	})
}
