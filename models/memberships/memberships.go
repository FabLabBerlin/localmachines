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
	AffectedMachines         string `orm:"type(text)"`
	AutoExtend               bool
	AutoExtendDurationMonths int64
	Archived                 bool
}

type RunningMembership struct {
	Id               int64
	Title            string
	StartDate        string
	TerminationDate  string
	UserId           int64
	UserLastName     string
	UserFirstName    string
	UserUsername     string
	UserEmail        string
	UserPhone        string
	UserCompany      string
}

func (this *Membership) AffectedMachineIdsLegacyDontUse() (ids []int64, err error) {
	parseErr := fmt.Errorf("cannot parse AffectedMachines property: '%v'", this.AffectedMachines)
	l := len(this.AffectedMachines)
	if l == 0 || this.AffectedMachines == "[]" {
		return []int64{}, nil
	}
	if l < 2 || this.AffectedMachines[0:1] != "[" || this.AffectedMachines[l-1:l] != "]" {
		return nil, parseErr
	}
	idStrings := strings.Split(this.AffectedMachines[1:l-1], ",")
	ids = make([]int64, 0, len(idStrings))
	for _, idString := range idStrings {
		if id, err := strconv.ParseInt(idString, 10, 64); err == nil {
			ids = append(ids, id)
		} else {
			return nil, fmt.Errorf("ParseInt: %v", err)
		}
	}
	return ids, nil
}

// Get an array of ID's of affected machines by the membership.
func (this *Membership) AffectedMachineIds() (ids []int64, err error) {

	ms, err := machine.GetAllAt(this.LocationId)
	if err != nil {
		return
	}

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

func (this *Membership) SetAffectedCategoryIds(ids []int64) error {
	strIds := make([]string, len(ids))
	for i, id := range ids {
		strIds[i] = strconv.FormatInt(id, 10)
	}
	this.AffectedCategories = "[" + strings.Join(strIds, ",") + "]"
	return nil
}

func (this *Membership) IsRndCentre() bool {
	return strings.Contains(this.Title, "R&D Centre")
}

// Returns whether the membership is affecting a machine
func (this *Membership) IsMachineAffected(m *machine.Machine) (bool, error) {
	if categoryIds, err := this.AffectedCategoryIds(); err == nil {
		for _, id := range categoryIds {
			if id == m.TypeId {
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

func init() {
	orm.RegisterModel(new(Membership))
}

func GetAllAt(locationId int64) (ms []*Membership, err error) {
	m := Membership{} // Used only for the TableName() method
	o := orm.NewOrm()
	_, err = o.QueryTable(m.TableName()).
		Filter("location_id", locationId).
		All(&ms)
	return
}

func GetAllRunningAt(locationId int64) (ms []*RunningMembership, err error) {
	// ms = []RunningMembership{}
	o := orm.NewOrm()

	query := "SELECT um.id, m.title, um.start_date, um.termination_date, "+
		"u.id AS user_id, u.last_name AS user_last_name, u.first_name AS user_first_name, u.username AS user_username, "+
		"u.email AS user_email, u.phone AS user_phone, u.company AS user_company "+
		"FROM user_memberships AS um "+
		"LEFT JOIN membership AS m ON um.membership_id = m.id "+
		"LEFT JOIN user AS u ON um.user_id = u.id "+
		"WHERE um.location_id=? AND (CAST(um.termination_date as date) >= NOW() OR um.termination_date IS NULL) "+
		"ORDER BY m.title, u.last_name"

	_, err = o.Raw(query, locationId).QueryRows(&ms)
	return
}

func Create(locId int64, name string) (m *Membership, err error) {
	m = &Membership{
		LocationId:               locId,
		Title:                    name,
		AutoExtend:               true,
		DurationMonths:           1,
		AutoExtendDurationMonths: 1,
	}

	_, err = orm.NewOrm().Insert(m)

	return
}

func Get(is int64) (*Membership, error) {
	m := Membership{}
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=?", m.TableName())
	o := orm.NewOrm()
	err := o.Raw(sql, is).QueryRow(&m)
	if err != nil {
		return nil, fmt.Errorf("Failed to get membership")
	}
	return &m, nil
}

func (membership *Membership) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(membership)
	return
}
