package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// User membership model that has a mapping in the database.
type UserMembership struct {
	Id           int64 `orm:"auto";"pk"`
	UserId       int64
	MembershipId int64
	StartDate    time.Time `orm:"type(datetime)"`
	EndDate      time.Time `orm:"type(datetime)"`
	AutoExtend   bool
}

// Returns the database table name of the UserMembership model.
func (this *UserMembership) TableName() string {
	return "user_membership"
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
	Title                 string
	ShortName             string
	DurationMonths        int
	Unit                  string
	MonthlyPrice          float32
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

// Creates user membership from user ID,
// membership ID and start time. Returns created user membership ID,
// and an error object.
func CreateUserMembership(
	userId int64,
	membershipId int64,
	startTime time.Time) (userMembershipId int64, err error) {

	userMembershipId = 0
	err = nil

	beego.Trace("membershipId:", membershipId)

	userMembership := UserMembership{}
	userMembership.UserId = userId
	userMembership.MembershipId = membershipId
	userMembership.StartDate = startTime

	// We need to get base membership data first to calc some stuff
	var baseMembership *Membership
	baseMembership, err = GetMembership(membershipId)
	if err != nil {
		beego.Error("Error getting membership:", err)

		userMembershipId = 0
		err = fmt.Errorf("Failed to get membership")
		return
	}

	// Set the auto extend flag of user membership to the one of base membership
	userMembership.AutoExtend = baseMembership.AutoExtend

	// Calculate end date by using base membership
	userMembership.EndDate = startTime.AddDate(
		0, int(baseMembership.DurationMonths), 0)

	// Store the user membership to database
	o := orm.NewOrm()
	userMembershipId, err = o.Insert(&userMembership)
	if err != nil {
		beego.Error("Error creating new user membership: ", err)

		userMembershipId = 0
		err = fmt.Errorf("Failed to create user membership")
		return
	}

	return
}

func GetAllUserMemberships() (userMemberships []*UserMembership, err error) {
	o := orm.NewOrm()
	um := new(UserMembership)
	_, err = o.QueryTable(um.TableName()).All(&userMemberships)
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
		beego.Error("Failed to read user membership:", err)
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
	sql := fmt.Sprintf("SELECT um.*, m.title, m.short_name, m.duration_months, "+
		"m.monthly_price, m.machine_price_deduction, m.affected_machines "+
		"FROM %s AS um "+
		"JOIN %s m ON um.membership_id=m.id "+
		"WHERE um.user_id=?",
		um.TableName(), m.TableName())

	var userMemberships []UserMembershipCombo
	_, err := o.Raw(sql, userId).QueryRows(&userMemberships)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user memberships")
	}

	userMembershipList := UserMembershipList{}
	userMembershipList.Data = userMemberships

	return &userMembershipList, nil
}

// Deletes user membership from the database by consuming
// user membership ID.
func DeleteUserMembership(userMembershipId int64) error {
	var num int64
	var err error
	o := orm.NewOrm()

	beego.Trace("models.DeleteUserMembership: userMembershipId:", userMembershipId)

	userMembership := &UserMembership{}
	userMembership.Id = userMembershipId

	num, err = o.QueryTable(userMembership.TableName()).
		Filter("id", userMembershipId).Delete()
	if err != nil {
		beego.Error("Failed to delete user membership:", err)
		return errors.New("Failed to delete user membership")
	}
	if num <= 0 {
		beego.Error("Nothing deleted")
		return fmt.Errorf("Nothing deleted")
	}

	beego.Trace("Deleted num rows:", num)
	return nil
}

// Updates user membership in the database by using a pointer
// to user membership store.
func UpdateUserMembership(userMembership *UserMembership) error {
	var err error
	var num int64

	o := orm.NewOrm()
	num, err = o.Update(userMembership)
	if err != nil {
		return err
	}

	beego.Trace("Rows affected:", num)
	return nil
}

// Automatically extend user membership end date if auto_extend for the specific
// membership is true and the end_date is before current date.
func AutoExtendUserMemberships() error {

	beego.Info("Running AutoExtendUserMemberships Task")
	currentTime := time.Now()

	// Get user memberships for all users
	um := UserMembership{} // Use this for the TableName() func only
	m := Membership{}      // Same goes for this one

	var userMembershipsArr []UserMembership
	o := orm.NewOrm()
	num, err := o.QueryTable(um.TableName()).All(&userMembershipsArr)
	if err != nil {
		beego.Error("Failed to get all user memberships:", err)
		return fmt.Errorf("Failde to get all user memberships")
	}

	beego.Trace("Total number of user memberships:", num)

	// Loop through the user memberships and check the end date
	for i := 0; i < len(userMembershipsArr); i++ {
		userMembership := userMembershipsArr[i]
		if userMembership.AutoExtend == true {

			// Check whether we need to extend the membership after all
			// by comparing membership end date with current date.
			if userMembership.EndDate.Before(currentTime) {
				beego.Trace("Extending user membership with ID:", userMembership.Id)
				beego.Trace("Base membership ID:", userMembership.MembershipId)
				beego.Trace("Current user membership end date:", userMembership.EndDate)

				// Get the amount of days we should extend from the base membership
				sql := fmt.Sprintf("SELECT auto_extend_duration_months FROM %s WHERE id=?",
					m.TableName())
				var autoExtendDuration int64
				beego.Trace("userMembership.MembershipId:", userMembership.MembershipId)
				err = o.Raw(sql, userMembership.MembershipId).QueryRow(&autoExtendDuration)
				if err != nil {
					beego.Error("Failed to exec raw query:", err)
					return fmt.Errorf("Failed to exec raw query")
				}
				beego.Trace("Auto extend duration in months:", autoExtendDuration)

				// Extend the user membership end date by number of days stored
				// in the base membership
				userMembership.EndDate = userMembership.EndDate.AddDate(0, int(autoExtendDuration), 0)
				beego.Trace("Extended user membership end date:", userMembership.EndDate)
				_, err = o.Update(&userMembership, "end_date")
				if err != nil {
					beego.Error("Failed to update user membership end date:", err)
					return fmt.Errorf("Failed to update user membership end date")
				}
				beego.Trace("Done extending user membership. Saved in the database.")
			}

		}
	}

	return nil
}
