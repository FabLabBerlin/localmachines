package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Membership struct {
	Id                       int64  `orm:"auto";"pk"`
	Title                    string `orm:"size(100)"`
	ShortName                string `orm:"size(100)"`
	DurationMonths           int64
	MonthlyPrice             float32
	MachinePriceDeduction    int
	AffectedMachines         string `orm:"type(text)"`
	AutoExtend               bool
	AutoExtendDurationMonths int64
}

func (this *Membership) AffectedMachineIds() ([]int64, error) {
	parseErr := fmt.Errorf("cannot parse AffectedMachines property: '%v'", this.AffectedMachines)
	l := len(this.AffectedMachines)
	if l < 2 || this.AffectedMachines[0:1] != "[" || this.AffectedMachines[l-1:l] != "]" {
		return nil, parseErr
	}
	idStrings := strings.Split(this.AffectedMachines[1:l-1], ",")
	ids := make([]int64, 0, len(idStrings))
	for _, idString := range idStrings {
		if id, err := strconv.ParseInt(idString, 10, 64); err == nil {
			ids = append(ids, id)
		} else {
			return nil, fmt.Errorf("ParseInt: %v", err)
		}
	}
	return ids, nil
}

func (this *Membership) IsMachineAffected(machineId int64) (bool, error) {
	if machineIds, err := this.AffectedMachineIds(); err == nil {
		for _, id := range machineIds {
			if id == machineId {
				return true, nil
			}
		}
		return false, nil
	} else {
		return false, err
	}
}

func (this *Membership) TableName() string {
	return "membership"
}

type UserMembership struct {
	Id           int64 `orm:"auto";"pk"`
	UserId       int64
	MembershipId int64
	StartDate    time.Time `orm:"type(datetime)"`
	EndDate      time.Time `orm:"type(datetime)"`
	AutoExtend   bool
}

func (this *UserMembership) TableName() string {
	return "user_membership"
}

func init() {
	orm.RegisterModel((new(Membership)), new(UserMembership))
}

// Extended UserMembership type that contains the same fields as UserMembership
// model, but additionaly fields from the membership model.
type UserMembershipCombo struct {
	Id                    int64
	UserId                int64
	MembershipId          int64
	StartDate             time.Time
	EndDate               time.Time
	Title                 string
	ShortName             string
	Duration              int
	Unit                  string
	MonthlyPrice          float32
	MachinePriceDeduction int
	AffectedMachines      string
	AutoExtend            bool
}

type UserMembershipList struct {
	Data []UserMembershipCombo
}

// Get all memberships from database
func GetAllMemberships() (memberships []*Membership, err error) {
	m := Membership{} // Used only for the TableName() method
	o := orm.NewOrm()
	num, err := o.QueryTable(m.TableName()).All(&memberships)
	if err != nil {
		beego.Error("Failed to get all memberships:", err)
		return nil, errors.New("Failed to get all memberships")
	}
	beego.Trace("Got num memberships:", num)
	return
}

// Creates a new membership in the database
func CreateMembership(membershipName string) (int64, error) {
	o := orm.NewOrm()
	membership := Membership{}
	membership.Title = membershipName
	membership.AutoExtend = true
	membership.DurationMonths = 1
	membership.AutoExtendDurationMonths = 1

	id, err := o.Insert(&membership)
	beego.Trace("Created membershipId:", id)
	if err == nil {
		return id, nil
	} else {
		return 0, fmt.Errorf("Failed to create membership: %v", err)
	}
}

// Get single membership from database using membership unique ID
func GetMembership(membershipId int64) (*Membership, error) {
	membership := Membership{}
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=?", membership.TableName())
	o := orm.NewOrm()
	err := o.Raw(sql, membershipId).QueryRow(&membership)
	if err != nil {
		beego.Error("Failed to get membership:", err, membershipId)
		return nil, fmt.Errorf("Failed to get membership")
	}
	return &membership, nil
}

// Updates membership in the database
func UpdateMembership(membership *Membership) error {
	var err error
	var num int64

	o := orm.NewOrm()
	num, err = o.Update(membership)
	if err != nil {
		return err
	}

	beego.Trace("Rows affected:", num)
	return nil
}

// Delete membership from the database
func DeleteMembership(membershipId int64) error {
	var num int64
	var err error
	o := orm.NewOrm()

	// Check if the membership has been added to user memberships.
	// Do not allow deletion if so.
	umem := UserMembership{}
	num, err = o.QueryTable(umem.TableName()).
		Filter("membership_id", membershipId).
		Count()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to get user memberships: %v", err))
	}

	if num > 0 {
		return errors.New(
			fmt.Sprintf("Membership has been assigned to %d users", num))
	}

	num, err = o.Delete(&Membership{Id: membershipId})
	if err != nil {
		return err
	}
	beego.Trace("Deleted num rows:", num)
	return nil
}

// Create user membership from an existing user membership struct.
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

// Get single user membership by using ID
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

// Get all user memberships for a user
func GetUserMemberships(userId int64) (*UserMembershipList, error) {

	o := orm.NewOrm()

	// Use these for the table names
	m := Membership{}
	um := UserMembership{}

	/*
		Id           int64 `orm:"auto";"pk"`
		UserId       int64
		MembershipId int64
		StartDate    time.Time
	*/

	// Joint query, select user memberships and expands them with
	// membership base information.
	sql := fmt.Sprintf("SELECT um.*, m.title, m.short_name, m.duration_months, "+
		"m.monthly_price, m.machine_price_deduction, m.affected_machines "+
		"FROM %s AS um "+
		"JOIN %s m ON um.membership_id=m.id "+
		"WHERE um.user_id=?",
		um.TableName(), m.TableName())

	var userMemberships []UserMembershipCombo
	var num int64
	var err error
	num, err = o.Raw(sql, userId).QueryRows(&userMemberships)
	if err != nil {
		beego.Error(err)
		return nil, fmt.Errorf("Failed to get user memberships")
	}
	beego.Trace("Got num user memberships:", num)

	// Loop through all of the user memberships and get the right time
	// by using raw queries. Convert start time and end time to UTC.
	type Dates struct {
		StartDate string
		EndDate   string
	}
	parseLayout := "2006-01-02 15:04:05"
	for i := 0; i < len(userMemberships); i++ {
		beego.Trace("---")
		beego.Trace(userMemberships[i].StartDate, userMemberships[i].EndDate)
		sql = fmt.Sprintf("SELECT start_date, end_date FROM %s WHERE id=?",
			um.TableName())
		var dates Dates
		err = o.Raw(sql, userMemberships[i].Id).QueryRow(&dates)
		if err != nil {
			beego.Error("Failed to execute raw query:", err)
		} else {
			userMemberships[i].StartDate, _ = time.ParseInLocation(parseLayout,
				dates.StartDate,
				time.UTC)
			userMemberships[i].EndDate, _ = time.ParseInLocation(parseLayout,
				dates.EndDate,
				time.UTC)
			beego.Trace(userMemberships[i].StartDate, userMemberships[i].EndDate)
		}
	}

	userMembershipList := UserMembershipList{}
	userMembershipList.Data = userMemberships

	return &userMembershipList, nil
}

// Delete membership from the database
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

// Updates membership in the database
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

	//sql := fmt.Sprintf(
	//	"SELECT id, membership_id, auto_extend FROM %s", um.TableName())

	var userMembershipsArr []UserMembership
	o := orm.NewOrm()
	num, err := o.QueryTable(um.TableName()).All(&userMembershipsArr)
	if err != nil {
		beego.Error("Failed to get all user memberships:", err)
		return fmt.Errorf("Failde to get all user memberships")
	}

	//num, err := o.Raw(sql).QueryRows(&userMembershipsArr)
	/*
		if err != nil {
			beego.Error("Failed to execute raw query:", err)
			return fmt.Errorf("Failed to exec raw query")
		}
	*/

	beego.Trace("Total number of user memberships:", num)

	// Loop through the user memberships and check the end date
	for i := 0; i < len(userMembershipsArr); i++ {
		userMembership := userMembershipsArr[i]
		if userMembership.AutoExtend == true {

			/*
				err = o.Read(&userMembership)
				if err != nil {
					beego.Trace("Failed to get end date of user membership:", err)
					return fmt.Errorf("Failed to get end date")
				}
			*/

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
