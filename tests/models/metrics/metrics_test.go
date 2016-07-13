package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestMetrics(t *testing.T) {

	Convey("Testing Metrics", t, func() {
		Reset(setup.ResetDB)

		Convey("Testing results", func() {
			m, err := machine.Create(1, "foo")
			if err != nil {
				panic(err.Error())
			}

			reservationPrice := 2.5
			m.ReservationPriceHourly = &reservationPrice
			if err := m.Update(false); err != nil {
				panic(err.Error())
			}

			nUsers := 10
			nPurchasesPerUser := 100

			for i := 0; i < nUsers; i++ {
				uid, err := users.CreateUser(&users.User{
					FirstName: "Joe",
					LastName:  "Doe",
					Email:     fmt.Sprintf("joe%v@example.com", i),
				})
				if err != nil {
					panic(err.Error())
				}

				_, err = user_locations.Create(&user_locations.UserLocation{
					LocationId: 1,
					UserId:     uid,
					UserRole:   "member",
				})
				if err != nil {
					panic(err.Error())
				}

				iv := &invutil.Invoice{}
				iv.LocationId = 1
				iv.UserId = uid
				iv.Month = 6
				iv.Year = 2016
				iv.Status = "outgoing"
				if _, err = invoices.Create(&iv.Invoice); err != nil {
					panic(err.Error())
				}

				for j := 0; j < nPurchasesPerUser; j++ {
					_, err := purchases.Create(&purchases.Purchase{
						LocationId:   1,
						Type:         purchases.TYPE_ACTIVATION,
						InvoiceId:    iv.Id,
						UserId:       uid,
						TimeStart:    time.Date(2016, 6, 1+j%20, 12, 13, 14, 0, time.UTC),
						TimeEnd:      time.Date(2016, 6, 1+j%20, 13, 13, 14, 0, time.UTC),
						Quantity:     1,
						PricePerUnit: 1,
						PriceUnit:    "hour",
					})
					if err != nil {
						panic(err.Error())
					}
				}
			}

			data, err := metrics.FetchData(1)
			if err != nil {
				panic(err.Error())
			}
			if err != nil {
				panic(err.Error())
			}

			resp, err := metrics.NewResponse(data)
			if err != nil {
				panic(err.Error())
			}
			if err != nil {
				panic(err.Error())
			}

			activationsJune, ok := resp.ActivationsByMonth["2016-06"]
			So(ok, ShouldBeTrue)
			So(activationsJune, ShouldEqual, nUsers*nPurchasesPerUser)
		})
	})
}
