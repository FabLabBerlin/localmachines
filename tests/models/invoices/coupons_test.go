package invoices

import (
	"fmt"
	"testing"
	"time"

	//"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/coupons"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/lib/fastbill/mock"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInvoiceCouponUsage(t *testing.T) {
	Convey("Test Generate and UseForInvoice", t, func() {
		Reset(setup.ResetDB)

		uid, err := users.CreateUser(&users.User{
			Email: "free@bars.com",
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
		cs, err := coupons.Generate(1, "foo", 10, 10)
		if err != nil {
			panic(err.Error())
		}
		c := cs[0]
		t := time.Now().AddDate(0, -1, 0)
		m := t.Month()
		y := t.Year()
		_, err = c.UseForInvoice(1.23, m, y)
		So(err, ShouldNotBeNil)
		if err = c.Assign(uid); err != nil {
			panic(err.Error())
		}
		usage, err := c.UseForInvoice(2.34, m, y)
		if err != nil {
			panic(err.Error())
		}
		So(usage.Value, ShouldEqual, 2.34)
	})

	Convey("Test invoice with coupon usage", t, func() {
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
		lasercutter, err := machine.Create(1, "Lasercutter")
		if err != nil {
			panic(err.Error())
		}
		lasercutterMinutes := 80
		lasercutterPricePerMinute := 1.0
		p := CreateTestPurchase(lasercutter.Id, "Lasercutter",
			time.Duration(lasercutterMinutes)*time.Minute,
			lasercutterPricePerMinute)
		p.UserId = uid
		o := orm.NewOrm()
		if _, err := o.Insert(p); err != nil {
			panic(err.Error())
		}
		cs, err := coupons.Generate(1, "foo", 10, 10)
		if err != nil {
			panic(err.Error())
		}
		c := cs[0]
		if err = c.Assign(uid); err != nil {
			panic(err.Error())
		}
		m := TIME_START.Month()
		y := TIME_START.Year()
		me := monthly_earning.MonthlyEarning{
			LocationId: 1,
			MonthFrom:  int(m),
			YearFrom:   y,
			MonthTo:    int(m),
			YearTo:     y,
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

		_, empty, err := monthly_earning.CreateFastbillDraft(invs[0])
		So(empty, ShouldBeFalse)
		So(err, ShouldBeNil)
		So(testServer.FbInv.Items, ShouldHaveLength, 1)
		rebate := 10 / (float64(lasercutterMinutes) * lasercutterPricePerMinute)
		So(testServer.FbInv.CashDiscountPercent,
			ShouldEqual,
			fmt.Sprintf("%v", 100*rebate))
	})
}
