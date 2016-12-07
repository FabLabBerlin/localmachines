package purchases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/email"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/FabLabBerlin/localmachines/models/users"
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

func (this *Reservation) Overlaps(other *Reservation) bool {
	if this.Purchase.MachineId != other.Purchase.MachineId {
		return false
	}

	if this.Purchase.Archived || this.Purchase.Cancelled || this.Purchase.ReservationDisabled ||
		other.Purchase.Archived || other.Purchase.Cancelled || other.Purchase.ReservationDisabled {

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

func (r *Reservation) SendEmailNotifications() (err error) {
	m, err := machine.Get(r.Purchase.MachineId)
	if err != nil {
		return fmt.Errorf("get machine: %v", err)
	}

	user, err := users.GetUser(r.Purchase.UserId)
	if err != nil {
		return fmt.Errorf("get user: %v", err)
	}

	location, err := locations.Get(r.Purchase.LocationId)
	if err != nil {
		return fmt.Errorf("get location: %v", err)
	}

	labSettings, err := settings.GetAllAt(r.Purchase.LocationId)
	if err != nil {
		return fmt.Errorf("get settings: %v", err)
	}

	timeLoc, err := time.LoadLocation(location.Timezone)
	if err != nil {
		timeLoc, err = time.LoadLocation("Europe/Berlin")
		if err != nil {
			err = nil
			timeLoc = time.UTC
		}
	}
	timeStart := r.Purchase.TimeStart.In(timeLoc).Format(time.RFC822)

	labEmail := email.New()
	labSubject := fmt.Sprintf("Reservation for %v on %v", m.Name, timeStart)
	labMessage := bytes.Buffer{}
	data := map[string]interface{}{
		"Location":    location,
		"User":        user,
		"Machine":     m,
		"Reservation": r,
		"TimeStart":   timeStart,
	}
	labTo := labSettings.GetString(r.Purchase.LocationId, settings.RESERVATION_NOTIFICATION_EMAIL)
	if labTo == nil || len(*labTo) < 4 {
		beego.Error("lab to mail address is nil")
		goto userNotification
	}
	err = email.LocationReservationNotification.Execute(&labMessage, data)
	if err != nil {
		beego.Error("execute location reservation notification template:", err)
		goto userNotification
	}

	err = labEmail.Send(*labTo, labSubject, labMessage.String())
	if err != nil {
		beego.Error("send location reservation notification:", err)
	}

userNotification:

	userMessage := bytes.Buffer{}
	userEmail := email.New()
	userSubject := fmt.Sprintf("%v reservation on %v",
		location.Title, timeStart)
	err = email.UserReservationNotification.Execute(&userMessage, data)
	if err != nil {
		beego.Error("execute user reservation notification template:", err)
		return err
	}

	err = userEmail.Send(user.Email, userSubject, userMessage.String())
	if err != nil {
		beego.Error("send location reservation notification:", err)
		return err
	}

	return err
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
				other.Purchase.TimeStart, other.Purchase.TimeEnd())
		}
	}

	return nil
}
