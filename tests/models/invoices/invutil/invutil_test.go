package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/inv_user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"strconv"
	"testing"
	"time"
)

const locId int64 = 1

func init() {
	setup.ConfigDB()
}

func TestInvutilInvoices(t *testing.T) {
	Convey("Testing invutil.Invoice model", t, func() {
		Reset(setup.ResetDB)
		Convey("AssureUserHasDraftFor creates invoice for month", func() {
			if lenInvoicesDB() > 0 || lenUserMembershipsDB() > 0 {
				panic("Expected clean state to test in.")
			}

			mNow := time.Now().Month()
			yNow := time.Now().Year()
			mLast := time.Now().AddDate(0, -1, -1).Month()
			yLast := time.Now().AddDate(0, -1, -1).Year()
			user, iv, _ := createInvoiceWithMembership(yLast, mLast, 1)
			iv.Invoice.Current = true
			if err := iv.Invoice.Save(); err != nil {
				panic(err.Error())
			}

			if l := lenInvoicesDB(); l != 1 {
				panic(strconv.Itoa(l))
			}
			if l := lenUserMembershipsDB(); l != 1 {
				panic(strconv.Itoa(l))
			}

			if err := invutil.AssureUserHasDraftFor(locId, user, yNow, mNow); err != nil {
				panic(err.Error())
			}

			existingIvs, err := invutil.GetAllAt(1)
			if err != nil {
				panic(err.Error())
			}

			data := invutil.NewPrefetchedData(1)
			if err := data.Prefetch(); err != nil {
				panic(err.Error())
			}

			for _, inv := range existingIvs {
				err := inv.InvoiceUserMemberships(data)
				So(err, ShouldBeNil)
			}

			existingUms, err := user_memberships.GetAllAt(1)
			if err != nil {
				panic(err.Error())
			}

			existingIums, err := inv_user_memberships.GetAllAt(1)
			if err != nil {
				panic(err.Error())
			}

			if l := len(existingIvs); l != 2 {
				panic(strconv.Itoa(l))
			}
			if l := len(existingUms); l != 1 {
				panic(strconv.Itoa(l))
			}
			if l := len(existingIums); l != 2 {
				panic(strconv.Itoa(l))
			}

			So(existingIvs[0].Month, ShouldEqual, mLast)
			So(existingIvs[0].Year, ShouldEqual, yLast)
			So(existingIvs[0].Current, ShouldBeFalse)

			So(existingIvs[1].Month, ShouldEqual, mNow)
			So(existingIvs[1].Year, ShouldEqual, yNow)
			So(existingIvs[1].Current, ShouldBeTrue)

			So(existingIums[0].UserId, ShouldEqual, user.Id)
			So(existingIums[0].StartDate.Month(), ShouldEqual, mLast)
			So(existingIums[1].UserId, ShouldEqual, user.Id)
			So(existingIums[1].StartDate.Month(), ShouldEqual, mLast)
		})

		Convey("Memberships in 1st month half affect 1st half", func() {
			testInvoiceWithMembershipAndTestPurchase(true)
		})

		Convey("Memberships in 1st month half don't affect 2nd half", func() {
			testInvoiceWithMembershipAndTestPurchase(false)
		})
	})
}

func testInvoiceWithMembershipAndTestPurchase(purchaseInsideMembershipInterval bool) {
	if lenInvoicesDB() > 0 || lenUserMembershipsDB() > 0 {
		panic("Expected clean state to test in.")
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err.Error())
	}

	mNow := time.Now().Month()
	yNow := time.Now().Year()
	mLast := time.Now().AddDate(0, -1, -1).Month()
	yLast := time.Now().AddDate(0, -1, -1).Year()
	user, iv, _ := createInvoiceWithMembership(yLast, mLast, 15)
	fmt.Printf("staarrrrt @ %v-%v-%v\n", yLast, mLast, 15)
	iv.Invoice.Current = true
	if err := iv.Invoice.Save(); err != nil {
		panic(err.Error())
	}

	var timeStart time.Time
	if purchaseInsideMembershipInterval {
		timeStart = time.Date(yNow, mNow, 2, 14, 10, 0, 0, loc)
	} else {
		timeStart = time.Date(yNow, mNow, 16, 14, 10, 0, 0, loc)
	}

	iv, err = invutil.GetDraft(1, user.Id, timeStart)
	if err != nil {
		panic(err.Error())
	}

	purchase := purchases.Purchase{
		LocationId:   1,
		Type:         purchases.TYPE_ACTIVATION,
		InvoiceId:    iv.Id,
		MachineId:    1,
		Created:      time.Now(),
		UserId:       user.Id,
		TimeStart:    timeStart,
		Quantity:     2,
		PricePerUnit: 23,
		PriceUnit:    "minute",
		Vat:          19,
	}

	if err = purchases.Create(&purchase); err != nil {
		panic(err.Error())
	}

	iv, err = invutil.GetDraft(1, user.Id, timeStart)
	if err != nil {
		panic(err.Error())
	}

	if l := len(iv.Purchases); l != 1 {
		panic(strconv.Itoa(l))
	}

	if l := len(iv.InvUserMemberships); l != 1 {
		panic(strconv.Itoa(l))
	}

	if err = iv.CalculateTotals(); err != nil {
		panic(err.Error())
	}

	So(math.Abs(iv.Sums.Purchases.Undiscounted-46) < 0.01, ShouldBeTrue)

	if purchaseInsideMembershipInterval {
		So(math.Abs(iv.Sums.Purchases.PriceInclVAT) < 0.01, ShouldBeTrue)
	} else {
		So(math.Abs(iv.Sums.Purchases.PriceInclVAT-46) < 0.01, ShouldBeTrue)
	}

}

func createInvoiceWithMembership(year int, month time.Month, dayStart int) (
	user *users.User,
	iv *invutil.Invoice,
	um *user_memberships.UserMembership) {
	o := orm.NewOrm()
	loc, _ := time.LoadLocation("Europe/Berlin")
	membershipStart := time.Date(year, month, dayStart, 14, 0, 0, 0, loc)

	user = &users.User{
		FirstName: "Amen",
		LastName:  "Hesus",
		Email:     "amen@example.com",
	}
	userId, err := users.CreateUser(user)
	if err != nil {
		panic(err.Error())
	}
	_, err = user_locations.Create(&user_locations.UserLocation{
		UserId:     user.Id,
		LocationId: 1,
	})
	if err != nil {
		panic(err.Error())
	}

	iv = &invutil.Invoice{}
	iv.LocationId = locId
	iv.UserId = userId
	iv.Month = int(month)
	iv.Year = year
	iv.Status = "draft"
	if _, err := invoices.Create(&iv.Invoice); err != nil {
		panic(err.Error())
	}

	m, err := memberships.Create(1, "Test Membership")
	if err != nil {
		panic(err.Error())
	}
	m.DurationMonths = 1
	m.MachinePriceDeduction = 100
	m.AutoExtend = false
	m.AutoExtendDurationMonths = 30
	m.AffectedMachines = "[1,2,3]"
	if err := m.Update(); err != nil {
		panic(err.Error())
	}

	um, err = user_memberships.Create(o, userId, m.Id, iv.Id, membershipStart)
	if err != nil {
		panic(err.Error())
	}
	return
}

func lenInvoicesDB() int {
	existing, err := invoices.GetAllInvoices(1)
	if err != nil {
		panic(err.Error())
	}
	return len(existing)
}

func lenUserMembershipsDB() int {
	ums, err := user_memberships.GetAllAt(1)
	if err != nil {
		panic(err.Error())
	}
	return len(ums)
}
