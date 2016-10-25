package user_memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/astaxie/beego/orm"
	"time"
)

// User membership model that has a mapping in the database.
type UserMembership struct {
	Id                    int64
	UserId                int64
	MembershipId          int64
	StartDate             time.Time `orm:"type(datetime)"`
	TerminationDate       time.Time `orm:"type(datetime)"`
	InitialDurationMonths int
	AutoExtend            bool

	Created time.Time
	Updated time.Time
}

/*func (this UserMembership) CloneOrm(o orm.Ormer, invId int64, invStatus string) error {
	var clone UserMembership
	clone = this
	clone.Id = 0
	clone.InvoiceId = invId
	clone.InvoiceStatus = invStatus
	_, err := o.Insert(&clone)
	return err
}*/

// Returns the database table name of the UserMembership model.
func (this *UserMembership) TableName() string {
	return "user_memberships"
}

func (this *UserMembership) TerminationDateDefined() bool {
	return this.TerminationDate.Unix() > 0
}

// Updates user membership in the database by using a pointer
// to user membership store.
func (this *UserMembership) Update() (err error) {
	o := orm.NewOrm()
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
	if this.AutoExtend {
		return true
	}
	return this.StartDate.AddDate(0, this.InitialDurationMonths, 0).After(t)
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
		AutoExtend:   m.AutoExtend,

		Created: time.Now(),
		Updated: time.Now(),
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

func GetAllAtList(locId int64) (list *List, err error) {
	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all at: %v", err)
	}
	msbyId := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		msbyId[m.Id] = m
	}

	ums, err := GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("GetAllAt: %v", err)
	}

	userMemberships := make([]*Combo, 0, len(ums))
	for _, um := range ums {
		c := Combo{}

		c.Id = um.Id
		c.UserId = um.UserId
		c.MembershipId = um.MembershipId
		c.StartDate = um.StartDate
		c.TerminationDate = um.TerminationDate
		c.AutoExtend = um.AutoExtend

		m, ok := msbyId[um.MembershipId]
		if !ok {
			return nil, fmt.Errorf("cannot find membership with id %v", um.MembershipId)
		}

		c.LocationId = m.LocationId
		c.Title = m.Title
		c.ShortName = m.ShortName
		c.DurationMonths = m.DurationMonths
		c.MonthlyPrice = m.MonthlyPrice
		c.MachinePriceDeduction = m.MachinePriceDeduction
		c.AffectedMachines = m.AffectedMachines

		userMemberships = append(userMemberships, &c)
	}

	list = &List{
		Data: userMemberships,
	}

	return list, nil
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
func GetForInvoice(invoiceId int64) (*List, error) {

	o := orm.NewOrm()

	// Use these for the table names
	m := memberships.Membership{}
	um := UserMembership{}

	// Joint query, select user memberships and expands them with
	// membership base information.
	sql := fmt.Sprintf("SELECT um.*, m.location_id, m.title, m.short_name, m.duration_months, "+
		"m.monthly_price, m.machine_price_deduction, m.affected_machines, m.auto_extend "+
		"FROM %s AS um "+
		"JOIN %s m ON um.membership_id=m.id "+
		"WHERE um.invoice_id=?",
		um.TableName(), m.TableName())

	var userMemberships []*Combo
	if _, err := o.Raw(sql, invoiceId).QueryRows(&userMemberships); err != nil {
		return nil, fmt.Errorf("query rows: %v", err)
	}

	list := List{
		Data: userMemberships,
	}

	return &list, nil
}
