package invoices

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestFastbillInvoiceActivation(t *testing.T) {
	Convey("Testing createFastbillDraft", t, func() {
		Reset(setup.ResetDB)
		p := CreateTestPurchase(22, "Lasercutter",
			time.Duration(12)*time.Minute, 0.5)

		ps := purchases.Purchases{
			Data: []*purchases.Purchase{
				p,
			},
		}

		t := time.Now()
		me := monthly_earning.MonthlyEarning{
			MonthFrom: int(t.Month()),
			YearFrom:  t.Year(),
			MonthTo:   int(t.Month()),
			YearTo:    t.Year(),
		}

		invs, err := me.GetInvoices(ps)
		if err != nil {
			panic(err.Error())
		}
		if len(invs) != 1 {
			panic("expected 1")
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
