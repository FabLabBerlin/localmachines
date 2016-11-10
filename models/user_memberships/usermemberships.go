package user_memberships

import (
	"fmt"
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
	StartDate             time.Time               `orm:"type(datetime)"`
	TerminationDate       time.Time               `orm:"type(datetime)"`
	InitialDurationMonths int
	DurationMonths        *int `orm:"-" json:",omitempty"`

	Created time.Time
	Updated time.Time
}

func (this *UserMembership) TableName() string {
	return TABLE_NAME
}

func (this *UserMembership) TerminationDateDefined() bool {
	return this.TerminationDate.Unix() > 0
}

func (this *UserMembership) Update(o orm.Ormer) (err error) {
	_, err = o.Update(this)
	return
}

func (this UserMembership) ActiveAt(t time.Time) bool {
	if t.Before(this.StartDate) {
		return false
	}

	if this.TerminationDateDefined() {
		return this.TerminationDate.After(t)
	} else {
		return true
	}
}

func init() {
	orm.RegisterModel(new(UserMembership))
}

// Creates user membership from user ID, membership ID and start time.
func Create(o orm.Ormer, userId, membershipId int64, start time.Time) (*UserMembership, error) {
	m, err := memberships.Get(membershipId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get membership: %v", err)
	}

	um := UserMembership{
		LocationId:            m.LocationId,
		UserId:                userId,
		MembershipId:          membershipId,
		StartDate:             start,
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
	if this.TerminationDateDefined() {
		ms := this.DurationMonthsUntil(this.TerminationDate)
		return &ms
	} else {
		return nil
	}
}

func (this *UserMembership) DurationMonthsUntil(end time.Time) (ms int) {
	t := this.StartDate

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
	if this.TerminationDateDefined() {
		ms := this.DurationMonthsUntil(this.TerminationDate)
		months = &ms
		ds := float64(this.TerminationDate.Sub(this.StartDate).Hours()) / 24
		if ds < 0 {
			ds = 0
		}
		days = &ds
	}

	return
}
