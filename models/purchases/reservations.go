package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Reservation struct {
	json.Marshaler
	json.Unmarshaler
	Purchase Purchase
}

func (this *Reservation) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

func (this *Reservation) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.Purchase)
}

func (this *Reservation) Id() int64 {
	return this.Purchase.Id
}

func (this *Reservation) LocationId() int64 {
	return this.Purchase.LocationId
}

func (this *Reservation) UserId() int64 {
	return this.Purchase.UserId
}

type ReservationCreatedResponse struct {
	Id int64
}

func GetReservation(id int64) (reservation *Reservation, err error) {
	reservation = &Reservation{}
	reservation.Purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&reservation.Purchase)

	return
}

func GetAllReservationsAt(locationId int64) (reservations []*Reservation, err error) {
	o := orm.NewOrm()
	r := new(Reservation)
	var purchases []*Purchase
	_, err = o.QueryTable(r.Purchase.TableName()).
		Filter("location_id", locationId).
		Filter("type", TYPE_RESERVATION).
		All(&purchases)
	if err != nil {
		return
	}
	reservations = make([]*Reservation, 0, len(purchases))
	for _, purchase := range purchases {
		reservation := &Reservation{
			Purchase: *purchase,
		}
		reservations = append(reservations, reservation)
	}
	return
}

func CreateReservation(reservation *Reservation) (int64, error) {

	// Get the reservation_price_hourly of the machine being reserved
	machine, err := machine.Get(reservation.Purchase.MachineId)
	if err != nil {
		return 0, fmt.Errorf("get machine: %v", err)
	}

	reservation.Purchase.Type = TYPE_RESERVATION
	reservation.Purchase.PricePerUnit = *machine.ReservationPriceHourly / 2
	reservation.Purchase.PriceUnit = "30 minutes"
	reservation.Purchase.Quantity = reservation.Purchase.quantityFromTimes()

	o := orm.NewOrm()
	return o.Insert(&reservation.Purchase)
}

func (reservation *Reservation) Update() (err error) {
	o := orm.NewOrm()
	reservation.Purchase.Quantity = reservation.Purchase.quantityFromTimes()
	_, err = o.Update(&reservation.Purchase)
	return
}

func DeleteReservation(id int64, isAdmin bool) (err error) {

	// Do not allow to delete reservations
	// if they are in the past
	// or they are happening today

	var reservation *Reservation
	reservation, err = GetReservation(id)
	if err != nil {
		beego.Error("Failed to get reservation")
		return fmt.Errorf("Failed to get reservation: %v", err)
	}

	timeNow := time.Now()

	// Check if past reservation
	if reservation.Purchase.TimeEnd.Before(timeNow) {
		beego.Error("Can not delete reservation from the past")
		return fmt.Errorf("Can not delete reservation from the past")
	}

	// Check if happening today
	if timeNow.Day() == reservation.Purchase.TimeStart.Day() &&
		timeNow.Month() == reservation.Purchase.TimeStart.Month() &&
		timeNow.Year() == reservation.Purchase.TimeStart.Year() &&
		!isAdmin {

		beego.Error("Can not delete a reservation happening today")
		return fmt.Errorf("Can not delete a reservation happening today")
	}

	// If we have not returned yet, then let's delete
	return Delete(id)
}
