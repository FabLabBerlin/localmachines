package invoices

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/lib/fastbill/mock"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestFastbillInvoiceActivation(t *testing.T) {
	Convey("Test Fastbill Invoice Activation", t, func() {
		Reset(setup.ResetDB)

		uid, err := users.CreateUser(&users.User{
			Email: "foo@bar.com",
		})
		if err != nil {
			panic(err.Error())
		}
		_, err = user_locations.Create(&user_locations.UserLocation{
			UserId:     uid,
			LocationId: 1,
		})
		if err != nil {
			panic(err.Error())
		}
		i3, err := machine.Create(1, "i3")
		if err != nil {
			panic(err.Error())
		}
		lasercutter, err := machine.Create(1, "Lasercutter")
		if err != nil {
			panic(err.Error())
		}
		p := CreateTestPurchase(lasercutter.Id, "Lasercutter", time.Duration(12)*time.Minute, 0.5)
		p.UserId = uid
		o := orm.NewOrm()
		if _, err := o.Insert(p); err != nil {
			panic(err.Error())
		}

		Convey("Testing createFastbillDraft", func() {

			t := time.Now()
			me := monthly_earning.MonthlyEarning{
				LocationId: 1,
				MonthFrom:  int(t.Month()),
				YearFrom:   t.Year(),
				MonthTo:    int(t.Month()),
				YearTo:     t.Year(),
			}

			invs, err := me.NewInvoices(19)
			if err != nil {
				panic(err.Error())
			}
			if n := len(invs); n != 1 {
				panic(fmt.Sprintf("expected 1 but got %v", n))
			}

			invs[0].User.ClientId = 1

			testServer := mock.NewServer()

			_, empty, err := monthly_earning.CreateFastbillDraft(&me, invs[0])
			So(empty, ShouldBeFalse)
			So(err, ShouldBeNil)
			So(testServer.FbInv.Items, ShouldHaveLength, 1)
		})

		Convey("Flatrate Memberships in draft leave no 0 price items", func() {
			p := CreateTestPurchase(lasercutter.Id, "Lasercutter", time.Duration(34)*time.Hour, 0.5)
			p.UserId = uid
			if _, err := o.Insert(p); err != nil {
				panic(err.Error())
			}
			p = CreateTestPurchase(i3.Id, "i3", time.Duration(100)*time.Nanosecond, 0.1)
			p.UserId = uid
			if _, err := o.Insert(p); err != nil {
				panic(err.Error())
			}

			ms, err := models.CreateMembership(1, "Full Flatrate")
			if err != nil {
				panic(err.Error())
			}
			ms.MonthlyPrice = 150
			ms.DurationMonths = 12
			ms.MachinePriceDeduction = 100
			if err = ms.Update(); err != nil {
				panic(err.Error())
			}
			ms.AffectedMachines = fmt.Sprintf("[%v]", lasercutter.Id)
			if err = ms.Update(); err != nil {
				panic(err.Error())
			}
			startTime := time.Now().AddDate(0, -2, 0)
			_, err = models.CreateUserMembership(uid, ms.Id, startTime)
			if err != nil {
				panic(err.Error())
			}
			t := time.Now()
			me := monthly_earning.MonthlyEarning{
				LocationId: 1,
				MonthFrom:  int(t.Month()),
				YearFrom:   t.Year(),
				MonthTo:    int(t.Month()),
				YearTo:     t.Year(),
			}

			invs, err := me.NewInvoices(19)
			if err != nil {
				panic(err.Error())
			}
			if n := len(invs); n != 1 {
				panic(fmt.Sprintf("expected 1 but got %v", n))
			}

			invs[0].User.ClientId = 1

			testServer := mock.NewServer()

			fastbill.API_URL = testServer.URL()

			_, empty, err := monthly_earning.CreateFastbillDraft(&me, invs[0])
			So(empty, ShouldBeFalse)
			So(err, ShouldBeNil)
			So(testServer.FbInv.Items, ShouldHaveLength, 1)
			item := testServer.FbInv.Items[0]
			So(item.Description, ShouldEqual, "Full Flatrate Membership (unit: month)")
		})
	})
}
