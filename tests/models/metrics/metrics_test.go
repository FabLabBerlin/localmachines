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
			s := Simulation{
				LocationId:        1,
				nUsers:            100,
				nPurchasesPerUser: 10,
				pricePerPurchase:  1,
			}
			simulate(s)

			// Locations mustn't interfer
			simulate(Simulation{
				LocationId:        2,
				nUsers:            123,
				nPurchasesPerUser: 10,
				pricePerPurchase:  1,
			})

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
			So(activationsJune, ShouldEqual, float64(s.nUsers*s.nPurchasesPerUser)*s.pricePerPurchase)
		})
	})
}

type Simulation struct {
	LocationId        int64
	nUsers            int
	nPurchasesPerUser int
	pricePerPurchase  float64
}

func simulate(s Simulation) {
	m, err := machine.Create(1, "foo")
	if err != nil {
		panic(err.Error())
	}

	reservationPrice := 2.5
	m.ReservationPriceHourly = &reservationPrice
	if err := m.Update(false); err != nil {
		panic(err.Error())
	}

	nUsers := s.nUsers
	nPurchasesPerUser := s.nPurchasesPerUser

	for i := 0; i < nUsers; i++ {
		uid, err := users.CreateUser(&users.User{
			FirstName: "Joe",
			LastName:  "Doe",
			Email:     fmt.Sprintf("joe%v-%v@example.com", i, s.LocationId),
		})
		if err != nil {
			panic(err.Error())
		}

		_, err = user_locations.Create(&user_locations.UserLocation{
			LocationId: s.LocationId,
			UserId:     uid,
			UserRole:   "member",
		})
		if err != nil {
			panic(err.Error())
		}

		iv := &invutil.Invoice{}
		iv.LocationId = s.LocationId
		iv.UserId = uid
		iv.Month = 6
		iv.Year = 2016
		iv.Status = "outgoing"
		if _, err = invoices.Create(&iv.Invoice); err != nil {
			panic(err.Error())
		}

		for j := 0; j < nPurchasesPerUser; j++ {
			_, err := purchases.Create(&purchases.Purchase{
				LocationId:   s.LocationId,
				Type:         purchases.TYPE_ACTIVATION,
				InvoiceId:    iv.Id,
				UserId:       uid,
				TimeStart:    time.Date(2016, 6, 1+j%20, 12, 13, 14, 0, time.UTC),
				TimeEnd:      time.Date(2016, 6, 1+j%20, 13, 13, 14, 0, time.UTC),
				Quantity:     1,
				PricePerUnit: s.pricePerPurchase,
				PriceUnit:    "hour",
			})
			if err != nil {
				panic(err.Error())
			}
		}
	}
}
