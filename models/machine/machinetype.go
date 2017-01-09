package machine

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

const (
	TYPE_TABLE_NAME = "machine_types"
)

type Type struct {
	Id         int64
	LocationId int64
	ShortName  string `orm:"size(20)"`
	Name       string `orm:"size(255)"`
	Archived   bool
}

func init() {
	orm.RegisterModel(new(Type))
}

func (t *Type) TableName() string {
	return TYPE_TABLE_NAME
}

func (t *Type) Create() (err error) {
	if t.Name == "" {
		return fmt.Errorf("No name")
	}
	if t.LocationId <= 0 {
		return fmt.Errorf("No loc id")
	}

	o := orm.NewOrm()
	_, err = o.Insert(t)
	if err != nil {
		return fmt.Errorf("insert: %v", err)
	}

	return
}

func GetType(id int64) (t *Type, err error) {
	o := orm.NewOrm()
	t = &Type{Id: id}
	err = o.Read(t)
	return
}

func GetAllTypes(locId int64) (ts []*Type, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TYPE_TABLE_NAME).
		Filter("location_id", locId).
		All(&ts)
	return
}

func (t *Type) Archive() (err error) {
	o := orm.NewOrm()
	t.Archived = true
	_, err = o.Update(t)
	return
}

func (t *Type) Unarchive() (err error) {
	o := orm.NewOrm()
	t.Archived = false
	_, err = o.Update(t)
	return
}

func (t *Type) Update() (err error) {
	if t.LocationId <= 0 {
		return fmt.Errorf("No loc id")
	}
	o := orm.NewOrm()
	_, err = o.Update(t)
	return
}
