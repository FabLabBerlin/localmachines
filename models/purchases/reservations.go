package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego/orm"
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

func (this *Reservation) Overlaps(other *Reservation) bool {
	if this.Purchase.MachineId != other.Purchase.MachineId {
		return false
	}

	return !(this.Purchase.TimeStart.Unix() >= other.Purchase.TimeEnd().Unix() ||
		other.Purchase.TimeStart.Unix() >= this.Purchase.TimeEnd().Unix())
}

func (this *Reservation) UserId() int64 {
	return this.Purchase.UserId
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

func CreateReservation(r *Reservation) (id int64, err error) {
	if err := r.assertOk(); err != nil {
		return 0, err
	}

	// Get the reservation_price_hourly of the machine being reserved
	m, err := machine.Get(r.Purchase.MachineId)
	if err != nil {
		return 0, fmt.Errorf("get machine: %v", err)
	}

	r.Purchase.Type = TYPE_RESERVATION
	r.Purchase.PricePerUnit = *m.ReservationPriceHourly / 2
	r.Purchase.PriceUnit = "30 minutes"

	err = Create(&r.Purchase)
	id = r.Purchase.Id

	return
}

func (r *Reservation) Update() (err error) {
	if err := r.assertOk(); err != nil {
		return err
	}

	return Update(&r.Purchase)
}

func (r *Reservation) assertOk() (err error) {
	if r.Purchase.LocationId <= 0 {
		return fmt.Errorf("invalid location id: %v", r.LocationId())
	}

	if mid := r.Purchase.MachineId; mid <= 0 {
		return fmt.Errorf("invalid machine id: %v", mid)
	}

	rs, err := GetAllReservationsAt(r.LocationId())
	if err != nil {
		return fmt.Errorf("get all reservations at %v: %v", r.LocationId(), err)
	}

	for _, other := range rs {
		if other.Purchase.Id == r.Purchase.Id {
			continue
		}

		if r.Overlaps(other) {
			return fmt.Errorf("overlapping with %v (%v - %v)", other.Id(),
				other.Purchase.TimeStart, other.Purchase.TimeEnd)
		}
	}

	return nil
}
