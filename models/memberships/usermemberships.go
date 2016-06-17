package memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/astaxie/beego/orm"
	"time"
)

// User membership model that has a mapping in the database.
type UserMembership struct {
	Id           int64
	UserId       int64
	MembershipId int64
	StartDate    time.Time `orm:"type(datetime)"`
	EndDate      time.Time `orm:"type(datetime)"`
	AutoExtend   bool

	InvoiceId     int64
	InvoiceStatus string
}

func (this UserMembership) Interval() lib.Interval {
	return lib.Interval{
		MonthFrom: int(this.StartDate.Month()),
		YearFrom:  this.StartDate.Year(),
		MonthTo:   int(this.EndDate.Month()),
		YearTo:    this.EndDate.Year(),
	}
}

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

// Extended user membership type that contains the same fields as
// the UserMembership model and fields of the Membership model.
type UserMembershipCombo struct {
	Id            int64
	UserId        int64
	MembershipId  int64
	StartDate     time.Time
	EndDate       time.Time
	AutoExtend    bool
	InvoiceId     int64
	InvoiceStatus string

	LocationId            int64
	Title                 string
	ShortName             string
	DurationMonths        int
	Unit                  string
	MonthlyPrice          float64
	MachinePriceDeduction int
	AffectedMachines      string
}

func (umc *UserMembershipCombo) Interval() lib.Interval {
	return umc.UserMembership().Interval()
}

func (umc *UserMembershipCombo) UserMembership() UserMembership {
	return UserMembership{
		Id:           umc.Id,
		UserId:       umc.UserId,
		MembershipId: umc.MembershipId,
		StartDate:    umc.StartDate,
		EndDate:      umc.EndDate,
		AutoExtend:   umc.AutoExtend,

		InvoiceId:     umc.InvoiceId,
		InvoiceStatus: umc.InvoiceStatus,
	}
}

// List container for the UserMembershipCombo type. Beego swagger
// has problems with interpretting plain arrays as documentable
// data type.
type UserMembershipList struct {
	Data []UserMembershipCombo
}

// Creates user membership from user ID, membership ID and start time.
func CreateUserMembership(userId, membershipId, invoiceId int64, start time.Time) (*UserMembership, error) {
	if invoiceId <= 0 {
		return nil, fmt.Errorf("need (valid) invoice id")
	}

	m, err := GetMembership(membershipId)
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

	o := orm.NewOrm()
	if um.Id, err = o.Insert(&um); err != nil {
		return nil, fmt.Errorf("insert: %v", err)
	}

	return &um, nil
}

func GetAllUserMembershipsAt(locationId int64) (userMemberships []*UserMembership, err error) {
	o := orm.NewOrm()
	m := new(Membership)
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
		QueryRows(&userMemberships)
	return
}

// Gets pointer to filled user membership store
// by using user membership ID.
func GetUserMembership(id int64) (*UserMembership, error) {
	userMembership := UserMembership{
		Id: id,
	}
	err := orm.NewOrm().Read(&userMembership)
	return &userMembership, err
}

// Gets all user memberships for a user by consuming user ID.
func GetUserMembershipsForInvoice(invoiceId int64) (*UserMembershipList, error) {

	o := orm.NewOrm()

	// Use these for the table names
	m := Membership{}
	um := UserMembership{}

	// Joint query, select user memberships and expands them with
	// membership base information.
	sql := fmt.Sprintf("SELECT um.*, m.location_id, m.title, m.short_name, m.duration_months, "+
		"m.monthly_price, m.machine_price_deduction, m.affected_machines, m.auto_extend "+
		"FROM %s AS um "+
		"JOIN %s m ON um.membership_id=m.id "+
		"WHERE um.invoice_id=?",
		um.TableName(), m.TableName())

	var userMemberships []UserMembershipCombo
	if _, err := o.Raw(sql, invoiceId).QueryRows(&userMemberships); err != nil {
		return nil, fmt.Errorf("query rows: %v", err)
	}

	userMembershipList := UserMembershipList{
		Data: userMemberships,
	}

	return &userMembershipList, nil
}
