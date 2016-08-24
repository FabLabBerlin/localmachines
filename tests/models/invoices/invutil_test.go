package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
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
		Convey("AssureUsersHaveInvoiceFor creates invoice for month", func() {
			if lenInvoicesDB() > 0 || lenUserMembershipsDB() > 0 {
				panic("Expected clean state to test in.")
			}

			mNow := time.Now().Month()
			yNow := time.Now().Year()
			mLast := time.Now().AddDate(0, -1, 0).Month()
			yLast := time.Now().AddDate(0, -1, 0).Year()
			user, _, _ := createInvoiceWithMembership(yLast, mLast, 1)

			if l := lenInvoicesDB(); l != 1 {
				panic(strconv.Itoa(l))
			}
			if l := lenUserMembershipsDB(); l != 1 {
				panic(strconv.Itoa(l))
			}

			if err := invutil.AssureUsersHaveInvoiceFor(locId, yNow, mNow); err != nil {
				panic(err.Error())
			}

			existingIvs, err := invoices.GetAllInvoices(1)
			if err != nil {
				panic(err.Error())
			}

			existingUms, err := user_memberships.GetAllAt(1)
			if err != nil {
				panic(err.Error())
			}

			if l := len(existingIvs); l != 2 {
				panic(strconv.Itoa(l))
			}
			if l := len(existingUms); l != 2 {
				panic(strconv.Itoa(l))
			}

			So(existingIvs[0].Month, ShouldEqual, mLast)
			So(existingIvs[0].Year, ShouldEqual, yLast)
			So(existingIvs[0].Current, ShouldBeFalse)

			So(existingIvs[1].Month, ShouldEqual, mNow)
			So(existingIvs[1].Year, ShouldEqual, yNow)
			So(existingIvs[1].Current, ShouldBeTrue)

			So(existingUms[0].UserId, ShouldEqual, user.Id)
			So(existingUms[0].StartDate.Month(), ShouldEqual, mLast)
			So(existingUms[1].UserId, ShouldEqual, user.Id)
			So(existingUms[1].StartDate.Month(), ShouldEqual, mLast)
		})

		Convey("Memberships in 1st month half don't affect 2nd half", func() {
			if lenInvoicesDB() > 0 || lenUserMembershipsDB() > 0 {
				panic("Expected clean state to test in.")
			}

			mNow := time.Now().Month()
			yNow := time.Now().Year()
			mLast := time.Now().AddDate(0, -1, 0).Month()
			yLast := time.Now().AddDate(0, -1, 0).Year()
			user, iv, userMembership := createInvoiceWithMembership(yLast, mLast, 15)

			So(userMembership.EndDate.Year(), ShouldEqual, yNow)
			So(userMembership.EndDate.Month(), ShouldEqual, mNow)
			So(userMembership.EndDate.Day(), ShouldEqual, 15)

			loc, _ := time.LoadLocation("Europe/Berlin")

			timeStart := time.Date(yNow, mNow, 16, 14, 10, 0, 0, loc)
			timeEnd := timeStart.Add(2 * time.Minute)

			iv, err := invutil.GetDraft(1, user.Id, timeStart)
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
				TimeEnd:      timeEnd,
				Quantity:     2,
				PricePerUnit: 23,
				PriceUnit:    "minute",
				Vat:          19,
			}

			_, err = purchases.Create(&purchase)
			if err != nil {
				panic(err.Error())
			}

			iv, err = invutil.GetDraft(1, user.Id, timeStart)
			if err != nil {
				panic(err.Error())
			}

			if l := len(iv.Purchases); l != 1 {
				panic(strconv.Itoa(l))
			}

			if l := len(iv.UserMemberships.Data); l != 1 {
				panic(strconv.Itoa(l))
			}

			if err = iv.CalculateTotals(); err != nil {
				panic(err.Error())
			}

			fmt.Printf("iv.Purchases.PriceInclVAT=%v\n", iv.Sums.Purchases.PriceInclVAT)
			fmt.Printf("iv.Purchases.Undiscounted=%v\n", iv.Sums.Purchases.Undiscounted)

			So(math.Abs(iv.Sums.Purchases.PriceInclVAT-46) < 0.01, ShouldBeTrue)
		})
	})
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
	m.AutoExtend = true
	m.AutoExtendDurationMonths = 30
	m.AffectedMachines = fmt.Sprintf("[1,2,3]")
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
