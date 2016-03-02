package invoices

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/invoices"
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

		inv := invoices.Invoice{}

		userSummaries, err := inv.GetUserSummaries(ps)
		if err != nil {
			panic(err.Error())
		}
		if len(*userSummaries) != 1 {
			panic("expected 1")
		}

		(*userSummaries)[0].User.ClientId = 1

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

		_, empty, err := invoices.CreateFastbillDraft((*userSummaries)[0])
		if empty {
			panic("empty")
		}
		if err != nil {
			panic(err.Error())
		}
	})
}
