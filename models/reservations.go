package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Reservation))
}

type Reservation struct {
	Id        int64 `orm:"auto";"pk"`
	MachineId int64
	UserId    int64
	TimeStart time.Time `orm:"type(datetime)"`
	TimeEnd   time.Time `orm:"type(datetime)"`
	Created   time.Time `orm:"type(datetime)"`
}

func (this *Reservation) TableName() string {
	return "reservations"
}

func GetReservation(id int64) (reservation *Reservation, err error) {
	err = orm.NewOrm().Read(reservation)
	return
}

func GetAllReservations() (reservations []*Reservation, err error) {
	o := orm.NewOrm()
	r := new(Reservation)
	_, err = o.QueryTable(r.TableName()).All(&reservations)
	return
}

func CreateReservation(reservation *Reservation) (int64, error) {
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
