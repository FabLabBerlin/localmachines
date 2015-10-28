package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Reservation struct {
	Id                   int64 `orm:"auto";"pk"`
	MachineId            int64
	UserId               int64
	TimeStart            time.Time `orm:"type(datetime)"`
	TimeEnd              time.Time `orm:"type(datetime)"`
	Created              time.Time `orm:"type(datetime)"`
	CurrentPrice         float64
	CurrentPriceCurrency string
	CurrentPriceUnit     string
}

func init() {
	orm.RegisterModel(new(Reservation))
}

type ReservationCreatedResponse struct {
	Id int64
}

func (this *Reservation) Slots() int {
	duration := this.TimeEnd.Unix() - this.TimeStart.Unix()
	return int(duration / 1800)
}

func (this *Reservation) TableName() string {
	return "reservations"
}

func (this *Reservation) PriceUnit() string {
	return "30 minute slot"
}

func GetReservation(id int64) (reservation *Reservation, err error) {
	r := Reservation{Id: id}
	reservation = &r
	o := orm.NewOrm()
	err = o.Read(reservation)
	return
}

func GetAllReservations() (reservations []*Reservation, err error) {
	o := orm.NewOrm()
	r := new(Reservation)
	_, err = o.QueryTable(r.TableName()).All(&reservations)
	return
}

func GetAllReservationsBetween(startTime, endTime time.Time) (reservations []*Reservation, err error) {
	allReservations, err := GetAllReservations()
	if err != nil {
		return
	}
	reservations = make([]*Reservation, 0, len(allReservations))
	for _, reservation := range allReservations {
		if startTime.Before(reservation.TimeStart) && reservation.TimeEnd.Before(endTime) {
			reservations = append(reservations, reservation)
		}
	}
	return
}

func CreateReservation(reservation *Reservation) (int64, error) {

	// Get the reservation_price_hourly of the machine being reserved
	machine := Machine{Id: reservation.MachineId}
	err, _ := machine.Read()
	if err != nil {
		beego.Error("Failed to read machine")
		return 0, fmt.Errorf("Failed to read machine")
	}

	reservation.CurrentPrice = *machine.ReservationPriceHourly
	reservation.CurrentPriceCurrency = "€"      // Reserved for possible future use
	reservation.CurrentPriceUnit = "30 minutes" // Reserved for possible future use

	o := orm.NewOrm()
	return o.Insert(reservation)
}

func UpdateReservation(reservation *Reservation) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(reservation)
	return
}

func DeleteReservation(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&Reservation{Id: id})
	return
}
