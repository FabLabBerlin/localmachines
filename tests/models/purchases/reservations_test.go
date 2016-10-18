package purchases

import (
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestReservations(t *testing.T) {
	Convey("Testing Reservation model", t, func() {
		Convey("Overlapping reservations raise an error", func() {
			m, err := machine.Create(1, "foo")
			if err != nil {
				panic(err.Error())
			}

			reservationPrice := 2.5
			m.ReservationPriceHourly = &reservationPrice
			if err := m.Update(false); err != nil {
				panic(err.Error())
			}

			_, err = purchases.CreateReservation(&purchases.Reservation{
				Purchase: purchases.Purchase{
					LocationId: 1,
					InvoiceId:  123,
					MachineId:  m.Id,
					TimeStart:  time.Date(2015, 10, 1, 14, 15, 0, 0, time.UTC),
					Quantity:   6,
					PriceUnit:  "30 minutes",
				},
			})
			if err != nil {
				panic(err.Error())
			}

			_, err = purchases.CreateReservation(&purchases.Reservation{
				Purchase: purchases.Purchase{
					LocationId: 1,
					InvoiceId:  123,
					MachineId:  m.Id,
					TimeStart:  time.Date(2015, 10, 1, 14, 16, 0, 0, time.UTC),
					Quantity:   2,
					PriceUnit:  "30 minutes",
				},
			})
			So(err, ShouldNotBeNil)
		})

		Convey("Non-overlapping reservations are fine", func() {
			m, err := machine.Create(1, "foo")
			if err != nil {
				panic(err.Error())
			}

			reservationPrice := 2.5
			m.ReservationPriceHourly = &reservationPrice
			if err := m.Update(false); err != nil {
				panic(err.Error())
			}

			_, err = purchases.CreateReservation(&purchases.Reservation{
				Purchase: purchases.Purchase{
					LocationId: 1,
					InvoiceId:  123,
					MachineId:  m.Id,
					TimeStart:  time.Date(2016, 10, 1, 14, 0, 0, 0, time.UTC),
					Quantity:   2,
					PriceUnit:  "30 minutes",
				},
			})
			if err != nil {
				panic(err.Error())
			}

			_, err = purchases.CreateReservation(&purchases.Reservation{
				Purchase: purchases.Purchase{
					LocationId: 1,
					InvoiceId:  123,
					MachineId:  m.Id,
					TimeStart:  time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
					Quantity:   4,
					PriceUnit:  "30 minutes",
				},
			})
			So(err, ShouldBeNil)
		})
	})
}
