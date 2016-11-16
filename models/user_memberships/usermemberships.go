package user_memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/astaxie/beego/orm"
	"time"
)

const TABLE_NAME = "user_memberships"

type UserMembership struct {
	Id                    int64
	LocationId            int64
	UserId                int64
	MembershipId          int64
	Membership            *memberships.Membership `orm:"-" json:",omitempty"`
	StartDate             string
	TerminationDate       *string
	InitialDurationMonths int
	DurationMonths        *int `orm:"-" json:",omitempty"`

	Created time.Time
	Updated time.Time
}

func (this *UserMembership) TableName() string {
	return TABLE_NAME
}

func (this *UserMembership) Update(o orm.Ormer) (err error) {
	_, err = o.Update(this)
	return
}

func (this UserMembership) ActiveAt(d day.Day) bool {
	t := time.Date(d.Year(), d.Month(), d.Day(), 11, 0, 0, 0, time.UTC)
	return this.ActiveAtTime(t)
}

func (this UserMembership) ActiveAtTime(t time.Time) bool {
	if this.StartDay().AfterTime(t) {
		return false
	}

	if td := this.TerminationDay(); td != nil {
		return !this.TerminationDay().BeforeTime(t)
	} else {
		return true
	}
}

func init() {
	orm.RegisterModel(new(UserMembership))
}

// Creates user membership from user ID, membership ID and start time.
func Create(o orm.Ormer, userId, membershipId int64, start day.Day) (*UserMembership, error) {
	m, err := memberships.Get(membershipId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get membership: %v", err)
	}

	um := UserMembership{
		LocationId:            m.LocationId,
		UserId:                userId,
		MembershipId:          membershipId,
		StartDate:             start.String(),
		InitialDurationMonths: int(m.DurationMonths),

		Created: time.Now(),
		Updated: time.Now(),
	}

	if um.Id, err = o.Insert(&um); err != nil {
		return nil, fmt.Errorf("insert: %v", err)
	}

	um.load()

	return &um, nil
}

func GetAllAt(locationId int64) (ums []*UserMembership, err error) {
	if _, err = orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		All(&ums); err != nil {

		return
	}

	for _, um := range ums {
		um.load()
	}

	return
}

func GetAllAtDeep(locId int64) (ums []*UserMembership, err error) {
	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all at: %v", err)
	}
	msbyId := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		msbyId[m.Id] = m
	}

	ums, err = GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("GetAllAt: %v", err)
	}

	for _, um := range ums {
		var ok bool

		um.Membership, ok = msbyId[um.MembershipId]
		if !ok {
			return nil, fmt.Errorf("cannot find membership with id %v", um.MembershipId)
		}
	}

	return
}

func GetAllOfDeep(locId, userId int64) (ums []*UserMembership, err error) {
	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all at: %v", err)
	}
	msbyId := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		msbyId[m.Id] = m
	}

	if _, err = orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locId).
		Filter("user_id", userId).
		All(&ums); err != nil {

		return
	}

	for _, um := range ums {
		var ok bool

		um.Membership, ok = msbyId[um.MembershipId]
		if !ok {
			return nil, fmt.Errorf("cannot find membership with id %v", um.MembershipId)
		}
		um.load()
	}

	return
}

func Get(id int64) (*UserMembership, error) {
	userMembership := UserMembership{
		Id: id,
	}

	if err := orm.NewOrm().Read(&userMembership); err != nil {
		return nil, err
	}
	userMembership.load()

	return &userMembership, nil
}

func Delete(o orm.Ormer, id int64) (err error) {
	_, err = o.Delete(&UserMembership{Id: id})
	return
}

func (this *UserMembership) durationMonths() *int {
	if td := this.TerminationDay(); td != nil {
		ms := this.DurationMonthsUntil(*td)
		return &ms
	} else {
		return nil
	}
}

func (this *UserMembership) DurationMonthsUntil(end day.Day) (ms int) {
	t := this.StartDay()

	for {
		t = t.AddDate(0, 1, 0)

		if t.Before(end.AddDate(0, 0, 2)) {
			ms++
		} else {
			break
		}
	}

	return ms
}

func (this *UserMembership) load() {
	this.DurationMonths = this.durationMonths()
}

// DurationModMonths is |TerminationDate - StartDate| % month. It is undefined
// when no TerminationDate is defined, thus returning nil, nil.
func (this *UserMembership) DurationModMonths() (months *int, days *float64) {
	if td := this.TerminationDay(); td != nil {
		ms := this.DurationMonthsUntil(*td)
		months = &ms
		ds := float64(td.Sub(this.StartDay()).Hours()) / 24
		if ds < 0 {
			ds = 0
		}
		days = &ds
	}

	return
}

func (this *UserMembership) StartDay() (d day.Day) {
	if d, err := day.NewString(this.StartDate); err == nil {
		return d
	} else {
		fmt.Printf("UserMembership#StartDay: day.NewString\n: %v", err)
	}
	return
}

func (this *UserMembership) TerminationDay() *day.Day {
	if td := this.TerminationDate; td != nil {
		d, err := day.NewString(*td)

		if err == nil {
			return &d
		} else {
			return nil
		}
	} else {
		return nil
	}
}
