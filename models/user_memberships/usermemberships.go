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

	if this.TerminationDateDefined() && this.TerminationDate.Before(t) {
		return false
	}

	return this.StartDate.AddDate(0, this.InitialDurationMonths, 0).After(t)
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

	return &um, nil
}

func GetAllAt(locationId int64) (ums []*UserMembership, err error) {
	orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		All(&ums)

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
	}

	return
}

func Get(id int64) (*UserMembership, error) {
	userMembership := UserMembership{
		Id: id,
	}
	err := orm.NewOrm().Read(&userMembership)
	return &userMembership, err
}

func Delete(o orm.Ormer, id int64) (err error) {
	_, err = o.Delete(&UserMembership{Id: id})
	return
}
