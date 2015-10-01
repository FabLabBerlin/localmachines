package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(ReservationRule))
}

type ReservationRule struct {
	Id        int64 `orm:"auto";"pk"`
	MachineId int64
	From      time.Time `orm:"type(datetime)"`
	To        time.Time `orm:"type(datetime)"`
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
}

func (this *ReservationRule) TableName() string {
	return "reservation_rules"
}

func GetReservationRule(id int64) (reservationRule *ReservationRule, err error) {
	err = orm.NewOrm().Read(reservationRule)
	return
}

func GetAllReservationRules() (reservationRules []ReservationRule, err error) {
	o := orm.NewOrm()
	r := new(Reservation)
	_, err = o.QueryTable(r.TableName()).All(&reservationRules)
	return
}

func CreateReservationRule(reservationRule *ReservationRule) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(reservationRule)
}

func UpdateReservationRule(reservationRule *ReservationRule) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(reservationRule)
	return
}

func DeleteReservationRule(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&ReservationRule{Id: id})
	return
}
