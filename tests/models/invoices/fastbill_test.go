package invoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib/fastbill"
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

func TestFastbillInvoiceActivation(t *testing.T) {
	Convey("Testing createFastbillDraft", t, func() {
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
		mid, _ := machine.CreateMachine(1, "Lasercutter")
		p := CreateTestPurchase(mid, "Lasercutter", time.Duration(12)*time.Minute, 0.5)
		p.UserId = uid
		o := orm.NewOrm()
		if _, err := o.Insert(p); err != nil {
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

		invs, err := me.NewInvoices()
		if err != nil {
			panic(err.Error())
		}
		if n := len(invs); n != 1 {
			panic(fmt.Sprintf("expected 1 but got %v", n))
		}

		invs[0].User.ClientId = 1

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			} else {
				resp = fastbill.InvoiceCreateResponse{}
			}

			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				panic(err.Error())
			}

			w.WriteHeader(http.StatusOK)
		}))

		fastbill.API_URL = testServer.URL

		_, empty, err := monthly_earning.CreateFastbillDraft(&me, invs[0])
		So(empty, ShouldBeFalse)
		So(err, ShouldBeNil)
	})
}
