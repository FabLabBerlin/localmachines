package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type SpaceReservation struct {
	json.Marshaler
	json.Unmarshaler
	purchase Purchase
}

func (this *SpaceReservation) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

func (this *SpaceReservation) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.purchase)
}

func CreateSpaceReservation(sr *SpaceReservation) (int64, error) {

	// Get the reservation_price_hourly of the machine being reserved
	machine := Machine{Id: sr.purchase.MachineId}
	err, _ := machine.Read()
	if err != nil {
		beego.Error("Failed to read machine")
		return 0, fmt.Errorf("Failed to read machine")
	}

	sr.purchase.Type = PURCHASE_TYPE_SPACE_RESERVATION
	//reservation.purchase.PricePerUnit = *machine.ReservationPriceHourly
	//reservation.purchase.PriceUnit = "30 minutes"

	o := orm.NewOrm()
	return o.Insert(&sr.purchase)
}
