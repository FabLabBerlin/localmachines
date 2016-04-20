package invoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

type MockServer struct {
	testServer *httptest.Server
	fbInv      fastbill.Invoice
}

func (s *MockServer) URL() string {
	return s.testServer.URL
}

func NewMockServer() (s *MockServer) {
	s = &MockServer{}
	s.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req fastbill.Request
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := dec.Decode(&req); err != nil {
			panic(err.Error())
		}
		var resp interface{}
		if req.SERVICE == fastbill.SERVICE_CUSTOMER_GET {
			resp = fastbill.CustomerGetResponse{
				RESPONSE: fastbill.CustomerList{
					Customers: []fastbill.Customer{
						fastbill.Customer{
							CUSTOMER_ID: "123",
						},
					},
				},
			}
		} else if req.SERVICE == fastbill.SERVICE_INVOICE_GET {
			resp = fastbill.InvoiceCreateResponse{}
		} else if req.SERVICE == fastbill.SERVICE_INVOICE_CREATE {
			buf, err := json.Marshal(req.DATA)
			if err != nil {
				panic(err.Error())
			}
			if err := json.Unmarshal(buf, &s.fbInv); err != nil {
				panic(err.Error())
			}
			resp = fastbill.InvoiceCreateResponse{}
		} else {
			panic("unknown service")
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			panic(err.Error())
		}

		w.WriteHeader(http.StatusOK)
	}))
	return
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

			testServer := NewMockServer()

			fastbill.API_URL = testServer.URL()

			_, empty, err := monthly_earning.CreateFastbillDraft(&me, invs[0])
			So(empty, ShouldBeFalse)
			So(err, ShouldBeNil)
			So(testServer.fbInv.Items, ShouldHaveLength, 1)
		})

		Convey("Flatrate Memberships in draft leave no 0 price items", func() {
			p := CreateTestPurchase(lasercutter.Id, "Lasercutter", time.Duration(34)*time.Hour, 0.5)
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

			testServer := NewMockServer()

			fastbill.API_URL = testServer.URL()

			_, empty, err := monthly_earning.CreateFastbillDraft(&me, invs[0])
			So(empty, ShouldBeFalse)
			So(err, ShouldBeNil)
			So(testServer.fbInv.Items, ShouldHaveLength, 1)
			item := testServer.fbInv.Items[0]
			So(item.Description, ShouldEqual, "Full Flatrate Membership")
		})
	})
}
