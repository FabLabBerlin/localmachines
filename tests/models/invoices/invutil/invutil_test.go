package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/inv_user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/models/mocks"
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
			fmt.Printf("existingIums[0].StartDate=%v\n", existingIums[0].StartDate)
			fmt.Printf("existingIums[0].StartDay()=%s\n", existingIums[0].StartDay())
			So(existingIums[0].StartDay().Month().Month(), ShouldEqual, mLast)
			So(existingIums[1].UserId, ShouldEqual, user.Id)
			So(existingIums[1].StartDay().Month().Month(), ShouldEqual, mLast)
		})

		Convey("Memberships in 1st month half affect 1st half", func() {
			testInvoiceWithMembershipAndTestPurchase(true)
		})

		Convey("Memberships in 1st month half don't affect 2nd half", func() {
			testInvoiceWithMembershipAndTestPurchase(false)
		})
	})

	Convey("InvoiceUserMemberships", t, func() {
		Reset(setup.ResetDB)
		Convey("(Inv)UserMemberships don't get unnecessarily duplicated", func() {
			inv := mocks.LoadInvoice(4165)

			So(len(inv.InvUserMemberships), ShouldEqual, 1)

			data := &invutil.PrefetchedData{
				LocationId: 1,
				UmbsByUid:  make(map[int64][]*user_memberships.UserMembership),
				IumbsByUid: make(map[int64][]*inv_user_memberships.InvoiceUserMembership),
			}

			data.UmbsByUid[19] = []*user_memberships.UserMembership{
				mocks.LoadUserMembership(14),
				mocks.LoadUserMembership(15),
			}

			for _, ium := range inv.InvUserMemberships {
				data.IumbsByUid[19] = append(data.IumbsByUid[19], ium)
			}

			if err := inv.InvoiceUserMemberships(data); err != nil {
				panic(err.Error())
			}

			So(len(inv.InvUserMemberships), ShouldEqual, 1)
		})

		Convey("Membership from 11/30-12/31 billed exactly once", func() {
			mt := MembershipIntervalTest{}

			mt.First.Month = 11
			mt.First.Year = 2015
			mt.First.Expect.LenInvUserMemberships = 1
			mt.First.Expect.Price = 0

			mt.Second.Month = 12
			mt.Second.Year = 2015
			mt.Second.Expect.LenInvUserMemberships = 1
			mt.Second.Expect.Price = 123.45

			mt.Membership.From = day.New(2015, 11, 30)
			mt.Membership.To = day.New(2015, 12, 31)
			mt.Membership.MonthlyPrice = 123.45

			mt.Run()
		})

		Convey("Membership from 11/09-12/15 billed exactly once", func() {
			mt := MembershipIntervalTest{}

			mt.First.Month = 11
			mt.First.Year = 2015
			mt.First.Expect.LenInvUserMemberships = 1
			mt.First.Expect.Price = 123.45

			mt.Second.Month = 12
			mt.Second.Year = 2015
			mt.Second.Expect.LenInvUserMemberships = 1
			mt.Second.Expect.Price = 0

			mt.Membership.From = day.New(2015, 11, 9)
			mt.Membership.To = day.New(2015, 12, 15)
			mt.Membership.MonthlyPrice = 123.45

			mt.Run()
		})

		Convey("Membership from 11/01-12/31 billed exactly twice", func() {
			mt := MembershipIntervalTest{}

			mt.First.Month = 11
			mt.First.Year = 2015
			mt.First.Expect.LenInvUserMemberships = 1
			mt.First.Expect.Price = 123.45

			mt.Second.Month = 12
			mt.Second.Year = 2015
			mt.Second.Expect.LenInvUserMemberships = 1
			mt.Second.Expect.Price = 123.45

			mt.Membership.From = day.New(2015, 11, 1)
			mt.Membership.To = day.New(2015, 12, 31)
			mt.Membership.MonthlyPrice = 123.45

			mt.Run()
		})

		Convey("Membership from 11/02-12/30 is less than 2 months => bill 1x", func() {
			mt := MembershipIntervalTest{}

			mt.First.Month = 11
			mt.First.Year = 2015
			mt.First.Expect.LenInvUserMemberships = 1
			mt.First.Expect.Price = 123.45

			mt.Second.Month = 12
			mt.Second.Year = 2015
			mt.Second.Expect.LenInvUserMemberships = 1
			mt.Second.Expect.Price = 0

			mt.Membership.From = day.New(2015, 11, 2)
			mt.Membership.To = day.New(2015, 12, 30)
			mt.Membership.MonthlyPrice = 123.45

			mt.Run()
		})

		Convey("Membership of 1 month (-1..+15d) gets billed exactly once", func() {
			for startOffset := 0; startOffset < 30; startOffset++ {
				for tolerance := -1; tolerance <= 15; tolerance++ {
					mt := MembershipIntervalTest{}

					mt.Membership.From = day.New(2015, 11, 1+startOffset)
					mt.Membership.To = mt.Membership.From.AddDate(0, 1, tolerance)
					mt.Membership.MonthlyPrice = 123.45

					mt.First.Month = 11
					mt.First.Year = 2015
					mt.First.Expect.LenInvUserMemberships = 1
					mt.First.Expect.Price = 123.45

					mt.Second.Month = 12
					mt.Second.Year = 2015
					if mt.Membership.To.Before(day.New(2015, 12, 1)) {
						mt.Second.Expect.LenInvUserMemberships = 0
					} else {
						mt.Second.Expect.LenInvUserMemberships = 1
					}
					mt.Second.Expect.Price = 0

					if mt.Membership.From.After(day.New(2015, 11, 28)) {
						mt.First.Expect.Price = 0
						mt.Second.Expect.Price = 123.45
					}

					mt.Run()
				}
			}
		})
	})

	Convey("userMembershipActiveHere", t, func() {
		Reset(setup.ResetDB)
		Convey("Membership from 11/01-12/31 is active in Nov", func() {
			td := "2015-12-31"
			um := &user_memberships.UserMembership{
				StartDate:       "2015-11-01",
				TerminationDate: &td,
			}

			inv := invutil.Invoice{}
			inv.Month = 11
			inv.Year = 2015

			So(inv.UserMembershipActiveHere(um), ShouldBeTrue)
		})

		Convey("Membership from 11/01-12/31 is active in Dec", func() {
			td := "2015-12-31"
			um := &user_memberships.UserMembership{
				StartDate:       "2015-11-01",
				TerminationDate: &td,
			}

			inv := invutil.Invoice{}
			inv.Month = 12
			inv.Year = 2015

			So(inv.UserMembershipActiveHere(um), ShouldBeTrue)
		})
	})

	Convey("userMembershipGetsBilledHere", t, func() {
		Reset(setup.ResetDB)
		Convey("Membership from 11/30-12/31 not billed in Nov", func() {
			td := "2015-12-31"
			um := &user_memberships.UserMembership{
				StartDate:       "2015-11-30",
				TerminationDate: &td,
			}

			inv := invutil.Invoice{}
			inv.Month = 11
			inv.Year = 2015

			So(inv.UserMembershipGetsBilledHere(um), ShouldBeFalse)
		})

		Convey("Membership from 11/30-12/31 is billed in Dec", func() {
			td := "2015-12-31"
			um := &user_memberships.UserMembership{
				StartDate:       "2015-11-30",
				TerminationDate: &td,
			}

			inv := invutil.Invoice{}
			inv.Month = 12
			inv.Year = 2015

			So(inv.UserMembershipGetsBilledHere(um), ShouldBeTrue)
		})

		Convey("Membership from 11/01-12/31 is billed in Nov", func() {
			td := "2015-12-31"
			um := &user_memberships.UserMembership{
				StartDate:       "2015-11-01",
				TerminationDate: &td,
			}

			inv := invutil.Invoice{}
			inv.Month = 11
			inv.Year = 2015

			So(inv.UserMembershipGetsBilledHere(um), ShouldBeTrue)
		})

		Convey("Membership from 11/01-12/31 is billed in Dec", func() {
			td := "2015-12-31"
			um := &user_memberships.UserMembership{
				StartDate:       "2015-11-01",
				TerminationDate: &td,
			}

			inv := invutil.Invoice{}
			inv.Month = 12
			inv.Year = 2015

			So(inv.UserMembershipGetsBilledHere(um), ShouldBeTrue)
		})
	})
}

// Test cases from https://goo.gl/QIJ194
type MembershipIntervalTest struct {
	First struct {
		Month  int
		Year   int
		Expect struct {
			LenInvUserMemberships int
			Price                 float64
		}
	}

	Second struct {
		Month  int
		Year   int
		Expect struct {
			LenInvUserMemberships int
			Price                 float64
		}
	}

	Membership struct {
		From         day.Day
		To           day.Day
		MonthlyPrice float64
	}
}

var uniqueId = 0

func (t *MembershipIntervalTest) Run() {
	uniqueId++
	n := fmt.Sprintf("%v", uniqueId)
	user := &users.User{
		FirstName: "Gerhard" + n,
		LastName:  "Schröder" + n,
		Email:     "gerhardt" + n + "@schroeder.net",
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

	m, err := memberships.Create(1, "Foo")
	if err != nil {
		panic(err.Error())
	}
	m.MonthlyPrice = t.Membership.MonthlyPrice
	if err := m.Update(); err != nil {
		panic(err.Error())
	}

	o := orm.NewOrm()
	startDate := t.Membership.From
	terminationDate := t.Membership.To

	um, err := user_memberships.Create(o, userId, m.Id, startDate)
	if err != nil {
		panic(err.Error())
	}

	terminationDateString := terminationDate.String()
	um.TerminationDate = &terminationDateString

	if err := um.Update(o); err != nil {
		panic(err.Error())
	}

	inv1stId, err := invoices.Create(&invoices.Invoice{
		LocationId: 1,
		Month:      t.First.Month,
		Year:       t.First.Year,
		UserId:     userId,
		Status:     "draft",
	})
	if err != nil {
		panic(err.Error())
	}

	inv2ndId, err := invoices.Create(&invoices.Invoice{
		LocationId: 1,
		Month:      t.Second.Month,
		Year:       t.Second.Year,
		UserId:     userId,
		Status:     "draft",
	})
	if err != nil {
		panic(err.Error())
	}

	inv1st, err := invutil.Get(inv1stId)
	if err != nil {
		panic(err.Error())
	}

	inv2nd, err := invutil.Get(inv2ndId)
	if err != nil {
		panic(err.Error())
	}

	data := invutil.NewPrefetchedData(1)
	if err := data.Prefetch(); err != nil {
		panic(err.Error())
	}

	if err := inv1st.InvoiceUserMemberships(data); err != nil {
		panic(err.Error())
	}

	if err := inv2nd.InvoiceUserMemberships(data); err != nil {
		panic(err.Error())
	}

	So(len(inv1st.InvUserMemberships), ShouldEqual, t.First.Expect.LenInvUserMemberships)
	So(inv1st.Sums.All.PriceInclVAT, ShouldEqual, t.First.Expect.Price)
	So(len(inv2nd.InvUserMemberships), ShouldEqual, t.Second.Expect.LenInvUserMemberships)
	So(inv2nd.Sums.All.PriceInclVAT, ShouldEqual, t.Second.Expect.Price)
}

func testInvoiceWithMembershipAndTestPurchase(purchaseInsideMembershipInterval bool) {
	if lenInvoicesDB() > 0 || lenUserMembershipsDB() > 0 {
		panic("Expected clean state to test in.")
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err.Error())
	}

	//membershipStart := time.Date(2015, time.June, 15, 0, 0, 0, 0, loc)
	//membershipEnd := time.Date(2015, time.July, 15, 0, 0, 0, 0, loc)

	/*mNow := time.Now().Month()
	yNow := time.Now().Year()
	mLast := time.Now().AddDate(0, -1, -1).Month()
	yLast := time.Now().AddDate(0, -1, -1).Year()*/
	user, iv, _ := createInvoiceWithMembership(2015, 6, 15)
	fmt.Printf("staarrrrt @ %v-%v-%v\n", 2015, 6, 15)
	iv.Invoice.Current = true
	if err := iv.Invoice.Save(); err != nil {
		panic(err.Error())
	}

	var timeStart time.Time
	if purchaseInsideMembershipInterval {
		timeStart = time.Date(2015, time.June, 18, 14, 10, 0, 0, loc)
	} else {
		timeStart = time.Date(2015, time.June, 2, 14, 10, 0, 0, loc)
	}

	/*iv, err = invutil.GetDraft(1, user.Id, timeStart)
	if err != nil {
		panic(err.Error())
	}*/
	/*iv = &invutil.Invoice{}
	iv.LocationId = locId
	iv.UserId = user.Id
	iv.Month = 7
	iv.Year = 2015
	iv.Status = "draft"
	if _, err := invoices.Create(&iv.Invoice); err != nil {
		panic(err.Error())
	}*/

	purchase := purchases.Purchase{
		LocationId:   1,
		Type:         purchases.TYPE_ACTIVATION,
		InvoiceId:    iv.Id,
		MachineId:    1,
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
	membershipStart := day.New(year, month, dayStart)

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

	um, err = user_memberships.Create(o, userId, m.Id, membershipStart)
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
