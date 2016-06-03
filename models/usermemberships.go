package models

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
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

	InvoiceId     uint64
	InvoiceStatus string
}

func (this *UserMembership) Interval() lib.Interval {
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
	Id                    int64
	UserId                int64
	MembershipId          int64
	StartDate             time.Time
	EndDate               time.Time
	LocationId            int64
	Title                 string
	ShortName             string
	DurationMonths        int
	Unit                  string
	MonthlyPrice          float64
	MachinePriceDeduction int
	AffectedMachines      string
	AutoExtend            bool
}

// List container for the UserMembershipCombo type. Beego swagger
// has problems with interpretting plain arrays as documentable
// data type.
type UserMembershipList struct {
	Data []UserMembershipCombo
}

// Creates user membership from user ID, membership ID and start time.
func CreateUserMembership(userId, membershipId int64, startTime time.Time) (id int64, err error) {
	baseMembership, err := GetMembership(membershipId)
	if err != nil {
		return 0, fmt.Errorf("Failed to get membership")
	}

	userMembership := UserMembership{
		UserId:       userId,
		MembershipId: membershipId,
		StartDate:    startTime,
		EndDate:      startTime.AddDate(0, int(baseMembership.DurationMonths), 0),
		AutoExtend:   baseMembership.AutoExtend,
	}

	o := orm.NewOrm()
	return o.Insert(&userMembership)
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
func GetUserMembership(userMembershipId int64) (*UserMembership, error) {
	userMembership := UserMembership{}
	userMembership.Id = userMembershipId
	o := orm.NewOrm()
	err := o.Read(&userMembership)
	if err != nil {
		return nil, fmt.Errorf("Failed to read user membership")
	}
	return &userMembership, nil
}

// Gets all user memberships for a user by consuming user ID.
func GetUserMemberships(userId int64) (*UserMembershipList, error) {

	o := orm.NewOrm()

	// Use these for the table names
	m := Membership{}
	um := UserMembership{}

	// Joint query, select user memberships and expands them with
	// membership base information.
	sql := fmt.Sprintf("SELECT um.*, m.location_id, m.title, m.short_name, m.duration_months, "+
		"m.monthly_price, m.machine_price_deduction, m.affected_machines "+
		"FROM %s AS um "+
		"JOIN %s m ON um.membership_id=m.id "+
		"WHERE um.user_id=?",
		um.TableName(), m.TableName())

	var userMemberships []UserMembershipCombo
	if _, err := o.Raw(sql, userId).QueryRows(&userMemberships); err != nil {
		return nil, fmt.Errorf("Failed to get user memberships")
	}

	userMembershipList := UserMembershipList{
		Data: userMemberships,
	}

	return &userMembershipList, nil
}

// Automatically extend user membership end date if auto_extend for the specific
// membership is true and the end_date is before current date.
func AutoExtendUserMemberships() (err error) {

	beego.Info("Running AutoExtendUserMemberships Task")

	if err = autoExtendUserMemberships(); err != nil {
		beego.Error("Failed to get all user memberships:", err)
	}

	return
}

func autoExtendUserMemberships() (err error) {
	ls, err := locations.GetAll()
	if err != nil {
		return fmt.Errorf("get all locations: %v", err)
	}

	for _, l := range ls {
		if err := extendUserMembershipsAt(l.Id); err != nil {
			return fmt.Errorf("extend userMemberships at %v: %v", l.Id, err)
		}
	}

	return
}

func extendUserMembershipsAt(locId int64) (err error) {
	ums, err := GetAllUserMembershipsAt(locId)
	if err != nil {
		return fmt.Errorf("get all user memberships: %v", err)
	}

	for _, um := range ums {
		if !um.AutoExtend || um.EndDate.After(time.Now()) {
			continue
		}

		beego.Trace("Extending user membership with Id", um.Id)

		m, err := GetMembership(um.MembershipId)
		if err != nil {
			return fmt.Errorf("get membership: %v", err)
		}

		um.EndDate = um.EndDate.AddDate(0, int(m.AutoExtendDurationMonths), 0)
		if err = um.Update(); err != nil {
			return fmt.Errorf("Failed to update user membership end date")
		}
	}

	return
}
