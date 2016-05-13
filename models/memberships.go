package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// Main membership type/struct that contains information
// about a mase membership and is mapper to a table in
// the database.
type Membership struct {
	Id                       int64
	LocationId               int64
	Title                    string `orm:"size(100)"`
	ShortName                string `orm:"size(100)"`
	DurationMonths           int64
	MonthlyPrice             float64
	MachinePriceDeduction    int
	AffectedMachines         string `orm:"type(text)"`
	AutoExtend               bool
	AutoExtendDurationMonths int64
	Archived                 bool
}

// Get an array of ID's of affected machines by the membership.
func (this *Membership) AffectedMachineIds() ([]int64, error) {
	parseErr := fmt.Errorf("cannot parse AffectedMachines property: '%v'", this.AffectedMachines)
	l := len(this.AffectedMachines)
	if l == 0 || this.AffectedMachines == "[]" {
		return []int64{}, nil
	}
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

// Returns whether the membership is affecting a machine
// with the given machine ID.
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

// Returns the database table name of the Membership model.
func (this *Membership) TableName() string {
	return "membership"
}

func init() {
	orm.RegisterModel(new(Membership))
}

// Gets all base memberships from database.
func GetAllMembershipsAt(locationId int64) (memberships []*Membership, err error) {
	m := Membership{} // Used only for the TableName() method
	o := orm.NewOrm()
	_, err = o.QueryTable(m.TableName()).
		Filter("location_id", locationId).
		All(&memberships)
	return
}

// Creates a new membership in the database with the given name.
func CreateMembership(locationId int64, name string) (m *Membership, err error) {
	m = &Membership{
		LocationId:               locationId,
		Title:                    name,
		AutoExtend:               true,
		DurationMonths:           1,
		AutoExtendDurationMonths: 1,
	}

	_, err = orm.NewOrm().Insert(m)

	return
}

// Gets single membership from database using membership unique ID
func GetMembership(membershipId int64) (*Membership, error) {
	membership := Membership{}
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=?", membership.TableName())
	o := orm.NewOrm()
	err := o.Raw(sql, membershipId).QueryRow(&membership)
	if err != nil {
		return nil, fmt.Errorf("Failed to get membership")
	}
	return &membership, nil
}

// Updates membership in the database by using a pointer to
// an existing membership store.
func (membership *Membership) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(membership)
	return
}

// Delete membership from the database by using a
// membership ID.
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
