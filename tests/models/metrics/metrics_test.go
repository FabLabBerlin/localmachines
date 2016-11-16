package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	"github.com/astaxie/beego/orm"
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
				nNormalUsers:      100,
				nFlatrateUsers:    11,
				nPurchasesPerUser: 10,
				pricePerPurchase:  2,
			}
			simulate(s)

			// Locations mustn't interfer
			simulate(Simulation{
				LocationId:        2,
				nNormalUsers:      123,
				nPurchasesPerUser: 10,
				pricePerPurchase:  1,
			})

			data, err := metrics.FetchData(1)
			if err != nil {
				panic(err.Error())
			}

			resp, err := metrics.NewResponse(data)
			if err != nil {
				panic(err.Error())
			}

			So(resp.MembershipsByMonth["2016-06"], ShouldEqual, s.nFlatrateUsers*17)
			So(resp.MembershipCountsByMonth["2016-06"], ShouldEqual, s.nFlatrateUsers)
			So(resp.ActivationsByMonth["2016-06"], ShouldEqual, float64(s.nNormalUsers*s.nPurchasesPerUser)*s.pricePerPurchase)
			So(resp.MinutesByMonth["2016-06"], ShouldEqual, float64((s.nNormalUsers+s.nFlatrateUsers)*s.nPurchasesPerUser)*60)
		})
	})
}

type Simulation struct {
	LocationId        int64
	nNormalUsers      int
	nFlatrateUsers    int
	nPurchasesPerUser int
	pricePerPurchase  float64
}

func simulate(s Simulation) {
	m, err := machine.Create(1, "foo")
	if err != nil {
		panic(err.Error())
	}

	mb, err := memberships.Create(s.LocationId, "bar")
	if err != nil {
		panic(err.Error())
	}
	mb.AffectedMachines = fmt.Sprintf("[%v]", m.Id)
	mb.MachinePriceDeduction = 100
	mb.DurationMonths = 1
	mb.MonthlyPrice = 17
	if err := mb.Update(); err != nil {
		panic(err.Error())
	}

	reservationPrice := 2.5
	m.ReservationPriceHourly = &reservationPrice
	if err := m.Update(false); err != nil {
		panic(err.Error())
	}

	nUsers := s.nNormalUsers + s.nFlatrateUsers
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
		iv.Status = "draft"
		if _, err = invoices.Create(&iv.Invoice); err != nil {
			panic(err.Error())
		}

		if i < s.nFlatrateUsers {
			_, err := user_memberships.Create(orm.NewOrm(), uid, mb.Id, day.New(2016, 6, 1))
			if err != nil {
				panic(err.Error())
			}
		}

		for j := 0; j < nPurchasesPerUser; j++ {
			if err := purchases.Create(&purchases.Purchase{
				LocationId:   s.LocationId,
				Type:         purchases.TYPE_ACTIVATION,
				InvoiceId:    iv.Id,
				MachineId:    m.Id,
				UserId:       uid,
				TimeStart:    time.Date(2016, 6, 1+j%20, 12, 13, 14, 0, time.UTC),
				Quantity:     1,
				PricePerUnit: s.pricePerPurchase,
				PriceUnit:    "hour",
			}); err != nil {
				panic(err.Error())
			}
		}
	}
}
