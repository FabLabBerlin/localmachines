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
					TimeEnd:    time.Date(2015, 10, 1, 14, 18, 0, 0, time.UTC),
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
					TimeEnd:    time.Date(2015, 10, 1, 14, 17, 0, 0, time.UTC),
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
					TimeStart:  time.Date(2016, 10, 1, 14, 15, 0, 0, time.UTC),
					TimeEnd:    time.Date(2016, 10, 1, 14, 18, 0, 0, time.UTC),
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
					TimeStart:  time.Date(2016, 10, 1, 14, 13, 0, 0, time.UTC),
					TimeEnd:    time.Date(2016, 10, 1, 14, 15, 0, 0, time.UTC),
				},
			})
			So(err, ShouldBeNil)
		})
	})
}
