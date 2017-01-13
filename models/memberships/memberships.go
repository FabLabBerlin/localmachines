package memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type Membership struct {
	Id                       int64
	LocationId               int64
	Title                    string `orm:"size(100)"`
	ShortName                string `orm:"size(100)"`
	DurationMonths           int64
	MonthlyPrice             float64
	MachinePriceDeduction    int
	AffectedCategories       string `orm:"type(text)"`
	AutoExtend               bool
	AutoExtendDurationMonths int64
	Archived                 bool
}

// Get an array of ID's of affected machines by the membership.
func (this *Membership) AffectedMachineIds() (ids []int64, err error) {

	ms, err := machine.GetAllAt(this.LocationId)
	if err != nil {
		return
	}

	fmt.Printf("AffectedMachineIds: ms=%v\n", ms)

	cids, err := this.AffectedCategoryIds()
	if err != nil {
		return
	}

	ids = make([]int64, 0, (len(ms)+len(cids))/2)
	for _, cid := range cids {
		for _, m := range ms {
			if m.TypeId == cid {
				ids = append(ids, m.Id)
			}
		}
	}

	return ids, nil
}

// Get an array of ID's of affected categories by the membership.
func (this *Membership) AffectedCategoryIds() ([]int64, error) {
	parseErr := fmt.Errorf("cannot parse AffectedCategories property: '%v'", this.AffectedCategories)
	l := len(this.AffectedCategories)
	if l == 0 || this.AffectedCategories == "[]" {
		return []int64{}, nil
	}
	if l < 2 || this.AffectedCategories[0:1] != "[" || this.AffectedCategories[l-1:l] != "]" {
		return nil, parseErr
	}
	idStrings := strings.Split(this.AffectedCategories[1:l-1], ",")
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

// Set affected categories IDs of the membership.
func (this *Membership) SetAffectedCategoryIds(categoryIds []int64) error {
	strIds := make([]string, len(categoryIds))
	for i, id := range categoryIds {
		strIds[i] = strconv.FormatInt(id, 10)
	}
	this.AffectedCategories = "[" + strings.Join(strIds, ",") + "]"
	return nil
}

func (this *Membership) IsRndCentre() bool {
	return strings.Contains(this.Title, "R&D Centre")
}

// Returns whether the membership is affecting a machine
// with the given machine ID.
func (this *Membership) IsMachineAffected(machineId int64) (bool, error) {
	if machineIds, err := this.AffectedMachineIds(); err == nil {
		fmt.Printf("IsMachineAffected: machineIds=%v\n", machineIds)
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
func GetAllAt(locationId int64) (memberships []*Membership, err error) {
	m := Membership{} // Used only for the TableName() method
	o := orm.NewOrm()
	_, err = o.QueryTable(m.TableName()).
		Filter("location_id", locationId).
		All(&memberships)
	return
}

// Creates a new membership in the database with the given name.
func Create(locationId int64, name string) (m *Membership, err error) {
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
func Get(membershipId int64) (*Membership, error) {
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
