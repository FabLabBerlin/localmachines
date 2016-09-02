package user_memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/astaxie/beego/orm"
	"time"
)

// User membership model that has a mapping in the database.
type UserMembership struct {
	Id           int64
	UserId       int64
	MembershipId int64
	Membership   *memberships.Membership `orm:"-"`
	StartDate    time.Time               `orm:"type(datetime)"`
	EndDate      time.Time               `orm:"type(datetime)"`
	AutoExtend   bool

	InvoiceId     int64
	InvoiceStatus string
}

func (this UserMembership) CloneOrm(o orm.Ormer, invId int64, invStatus string) error {
	var clone UserMembership
	clone = this
	clone.Id = 0
	clone.InvoiceId = invId
	clone.InvoiceStatus = invStatus
	_, err := o.Insert(&clone)
	return err
}

/*func (this UserMembership) Interval() lib.Interval {
	return lib.Interval{
		MonthFrom: int(this.StartDate.Month()),
		YearFrom:  this.StartDate.Year(),
		MonthTo:   int(this.EndDate.Month()),
		YearTo:    this.EndDate.Year(),
	}
}*/

// Returns the database table name of the UserMembership model.
func (this *UserMembership) TableName() string {
	return "user_membership"
}

// Updates user membership in the database by using a pointer
// to user membership store.
func (this *UserMembership) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(this)
	return
}

func init() {
	orm.RegisterModel(new(UserMembership))
}

// Creates user membership from user ID, membership ID and start time.
func Create(o orm.Ormer, userId, membershipId, invoiceId int64, start time.Time) (*UserMembership, error) {
	if invoiceId <= 0 {
		return nil, fmt.Errorf("need (valid) invoice id")
	}

	m, err := memberships.Get(membershipId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get membership: %v", err)
	}

	um := UserMembership{
		UserId:       userId,
		MembershipId: membershipId,
		StartDate:    start,
		EndDate:      start.AddDate(0, int(m.DurationMonths), 0),
		AutoExtend:   m.AutoExtend,
		InvoiceId:    invoiceId,
	}

	if um.Id, err = o.Insert(&um); err != nil {
		return nil, fmt.Errorf("insert: %v", err)
	}

	return &um, nil
}

func GetAllAt(locationId int64) (ums []*UserMembership, err error) {
	o := orm.NewOrm()
	m := new(memberships.Membership)
	um := new(UserMembership)

	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, fmt.Errorf("new query builder: %v", err)
	}

	qb.Select(um.TableName() + ".*").
		From(um.TableName()).
		InnerJoin(m.TableName()).
		On("membership.id = membership_id").
		Where("location_id = ?")

	_, err = o.Raw(qb.String(), locationId).
		QueryRows(&ums)
	return
}

func GetAllAtDeepQuery(locId int64) (ums []*UserMembership, err error) {
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
		m, ok := msbyId[um.MembershipId]
		if !ok {
			return nil, fmt.Errorf("cannot find membership with id %v", um.MembershipId)
		}

		um.Membership = m
	}

	return
}

// Gets pointer to filled user membership store
// by using user membership ID.
func Get(id int64) (*UserMembership, error) {
	userMembership := UserMembership{
		Id: id,
	}
	err := orm.NewOrm().Read(&userMembership)
	return &userMembership, err
}

func Delete(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&UserMembership{Id: id})
	return
}

// Gets all user memberships for a user by consuming user ID.
func GetForInvoice(locId, invoiceId int64) (ums []*UserMembership, err error) {

	o := orm.NewOrm()

	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all at: %v", err)
	}
	msbyId := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		msbyId[m.Id] = m
	}

	// Use these for the table names
	um := UserMembership{}

	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, fmt.Errorf("new query builder: %v", err)
	}

	qb.Select(um.TableName() + ".*").
		From(um.TableName()).
		Where("invoice_id = ?")

	_, err = o.Raw(qb.String(), invoiceId).
		QueryRows(&ums)

	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}

	for _, um := range ums {
		m, ok := msbyId[um.MembershipId]
		if !ok {
			return nil, fmt.Errorf("cannot find membership with id %v", um.MembershipId)
		}

		um.Membership = m
	}

	return
}
