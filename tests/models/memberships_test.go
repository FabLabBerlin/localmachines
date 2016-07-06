package modelTest

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

// TODO: The way go convey tests are supposed to be written is more
// human readable than this. Improve on that.

func CreateMembershipsActivation(userId, machineId, invoiceId int64, startTime time.Time, minutes float64) (id int64) {

	activation := purchases.Activation{
		Purchase: purchases.Purchase{
			LocationId: 1,
			TimeStart:  startTime,
			UserId:     userId,
			MachineId:  machineId,
			InvoiceId:  invoiceId,
		},
	}
	activation.Purchase.TimeEnd = activation.Purchase.TimeStart.Add(time.Duration(minutes) * time.Minute)

	id, err := purchases.Create(&activation.Purchase)
	if err != nil {
		panic(err.Error())
	}
	return
}

func TestMemberships(t *testing.T) {
	Convey("Testing Membership model", t, func() {

		Reset(setup.ResetDB)

		Convey("Testing CreateMembership", func() {
			name := "Membership X"

			Convey("When creating a single membership", func() {
				membership, err := memberships.Create(1, name)
				if err != nil {
					panic(err.Error())
				}

				Convey("There should be no errors and the ID should be valid", func() {
					So(err, ShouldBeNil)
					So(membership.Id, ShouldBeGreaterThan, 0)
				})

				Convey("When reading it back by using the ID", func() {
					membership, err := memberships.Get(membership.Id)
					if err != nil {
						panic(fmt.Sprintf("%v ... membershipId: %v", err.Error(), membership.Id))
					}
					Convey("It should return no error", func() {
						So(err, ShouldBeNil)
					})

					Convey("It should return the membership", func() {
						So(membership, ShouldNotBeNil)

						var membershipType *memberships.Membership
						So(membership, ShouldHaveSameTypeAs, membershipType)
					})

					Convey("Title should equal the initially given one", func() {
						title := membership.Title
						So(title, ShouldEqual, name)
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

			name := "The Membership"

			Convey("Getting all memberships with empty database", func() {
				ms, err := memberships.GetAllAt(1)

				Convey("Should return no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should return an empty array", func() {
					So(len(ms), ShouldEqual, 0)
				})
			})

			Convey("Getting existing memberships", func() {
				memberships.Create(1, name)
				memberships.Create(1, name)
				ms, err := memberships.GetAllAt(1)

				Convey("Shoud return no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should return an array with 2 memberships", func() {
					So(len(ms), ShouldEqual, 2)
				})
			})
		})

		Convey("Testing GetMembership", func() {
			name := "The Membership"

			Convey("Getting non-existing membership", func() {
				_, err := memberships.Get(0)

				Convey("Should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Creating a membership and getting it", func() {
				m, _ := memberships.Create(1, name)
				membership, err := memberships.Get(m.Id)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Got membership ID should match the ID of the previously created", func() {
					So(membership.Id, ShouldEqual, m.Id)
				})
			})
		})

		Convey("Testing UpdateMembership", func() {
			name := "Update Membership"
			newName := "New Membership Name"

			Convey("Try to update non existing membership", func() {
				m := &memberships.Membership{
					Title: name,
				}
				err := m.Update()

				Convey("There should be error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Create membership and update it", func() {
				m, _ := memberships.Create(1, name)
				m.Title = newName
				err := m.Update()
				nm, _ := memberships.Get(m.Id)

				Convey("There should be no error", func() {
					So(err, ShouldBeNil)
				})

				Convey("The title of the read-back membership should equal given one", func() {
					So(nm.Title, ShouldEqual, newName)
				})
			})
		})

	})
}
