package coupons

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/coupons"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/lib/fastbill/mock"
	"github.com/FabLabBerlin/localmachines/tests/models/util"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

var TIME_START = util.TIME_START

func init() {
	setup.ConfigDB()
}

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
		inv := &invutil.Invoice{}
		inv.LocationId = 1
		inv.UserId = uid
		inv.Month = int(TIME_START.Month())
		inv.Year = TIME_START.Year()
		inv.Status = "outgoing"
		if _, err = invoices.Create(&inv.Invoice); err != nil {
			panic(err.Error())
		}

		p := util.CreateTestPurchase(lasercutter.Id, "Lasercutter",
			time.Duration(lasercutterMinutes)*time.Minute,
			lasercutterPricePerMinute)
		p.InvoiceId = inv.Id
		p.UserId = uid

		if err := purchases.Create(p); err != nil {
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
		invs, err := invutil.GetAllOfMonthAt(1, y, m)
		if err != nil {
			panic(err.Error())
		}
		if n := len(invs); n != 1 {
			panic(fmt.Sprintf("expected 1 but got %v", n))
		}
		invs[0].User.ClientId = 1

		testServer := mock.NewServer("foo@bar.com")

		_, empty, err := invs[0].FastbillCreateDraft(false)
		So(empty, ShouldBeFalse)
		So(err, ShouldBeNil)
		So(testServer.FbInv.Items, ShouldHaveLength, 1)
		rebate := 10 / (float64(lasercutterMinutes) * lasercutterPricePerMinute)
		So(testServer.FbInv.CashDiscountPercent,
			ShouldEqual,
			fmt.Sprintf("%v", 100*rebate))
	})
}
