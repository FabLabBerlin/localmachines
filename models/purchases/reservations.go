package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Reservation struct {
	json.Marshaler
	json.Unmarshaler
	purchase Purchase
}

func (this *Reservation) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

func (this *Reservation) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.purchase)
}

func (this *Reservation) LocationId() int64 {
	return this.purchase.LocationId
}

func (this *Reservation) UserId() int64 {
	return this.purchase.UserId
}

type ReservationCreatedResponse struct {
	Id int64
}

func GetReservation(id int64) (reservation *Reservation, err error) {
	reservation = &Reservation{}
	reservation.purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&reservation.purchase)

	return
}

func GetAllReservationsAt(locationId int64) (reservations []*Reservation, err error) {
	o := orm.NewOrm()
	r := new(Reservation)
	var purchases []*Purchase
	_, err = o.QueryTable(r.purchase.TableName()).
		Filter("location_id", locationId).
		Filter("type", TYPE_RESERVATION).
		All(&purchases)
	if err != nil {
		return
	}
	reservations = make([]*Reservation, 0, len(purchases))
	for _, purchase := range purchases {
		reservation := &Reservation{
			purchase: *purchase,
		}
		reservations = append(reservations, reservation)
	}
	return
}

/*func GetAllReservationsBetween(startTime, endTime time.Time) (reservations []*Reservation, err error) {
	allReservations, err := GetAllReservations()
	if err != nil {
		return
	}
	reservations = make([]*Reservation, 0, len(allReservations))
	for _, reservation := range allReservations {
		if startTime.Before(reservation.purchase.TimeStart) && reservation.purchase.TimeEnd.Before(endTime) {
			reservations = append(reservations, reservation)
		}
	}
	return
}*/

func CreateReservation(reservation *Reservation) (int64, error) {

	// Get the reservation_price_hourly of the machine being reserved
	machine := models.Machine{Id: reservation.purchase.MachineId}
	err, _ := machine.Read()
	if err != nil {
		beego.Error("Failed to read machine")
		return 0, fmt.Errorf("Failed to read machine")
	}

	reservation.purchase.Type = TYPE_RESERVATION
	reservation.purchase.PricePerUnit = *machine.ReservationPriceHourly / 2
	reservation.purchase.PriceUnit = "30 minutes"
	reservation.purchase.Quantity = reservation.purchase.quantityFromTimes()

	o := orm.NewOrm()
	return o.Insert(&reservation.purchase)
}

func (reservation *Reservation) Update() (err error) {
	o := orm.NewOrm()
	reservation.purchase.Quantity = reservation.purchase.quantityFromTimes()
	_, err = o.Update(&reservation.purchase)
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
	if reservation.purchase.TimeEnd.Before(timeNow) {
		beego.Error("Can not delete reservation from the past")
		return fmt.Errorf("Can not delete reservation from the past")
	}

	// Check if happening today
	if timeNow.Day() == reservation.purchase.TimeStart.Day() &&
		timeNow.Month() == reservation.purchase.TimeStart.Month() &&
		timeNow.Year() == reservation.purchase.TimeStart.Year() &&
		!isAdmin {

		beego.Error("Can not delete a reservation happening today")
		return fmt.Errorf("Can not delete a reservation happening today")
	}

	// If we have not returned yet, then let's delete
	return Delete(id)
}
