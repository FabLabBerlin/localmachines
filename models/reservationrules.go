package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(ReservationRule))
}

/*
 * ReservationRule
 *
 * DateStart, DateEnd strings of "YYYY-MM-DD"
 * TimeStart, TimeEnd strings of "HH:MM"
 *
 */
type ReservationRule struct {
	Id          int64
	LocationId  int64
	Name        string
	MachineId   int64
	Available   bool
	Unavailable bool
	DateStart   string
	DateEnd     string
	TimeStart   string
	TimeEnd     string
	TimeZone    string
	Monday      bool
	Tuesday     bool
	Wednesday   bool
	Thursday    bool
	Friday      bool
	Saturday    bool
	Sunday      bool
	Created     time.Time `orm:"type(datetime)"`
}

func (this *ReservationRule) TableName() string {
	return "reservation_rules"
}

func GetReservationRule(id int64) (*ReservationRule, error) {
	rule := ReservationRule{Id: id}
	err := orm.NewOrm().Read(&rule)
	return &rule, err
}

func GetAllReservationRulesAt(locId int64) (rules []ReservationRule, err error) {
	o := orm.NewOrm()
	r := new(ReservationRule)
	_, err = o.QueryTable(r.TableName()).
		Filter("location_id", locId).
		All(&rules)
	return
}

func CreateReservationRule(reservationRule *ReservationRule) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(reservationRule)
}

func (this *ReservationRule) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(this)
	return
}

func DeleteReservationRule(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&ReservationRule{Id: id})
	return
}

type ReservationRuleCreatedResponse struct {
	Id int64
}
