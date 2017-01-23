package categories

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

const (
	TYPE_TABLE_NAME = "categories"
)

type Category struct {
	Id         int64
	LocationId int64
	ShortName  string `orm:"size(20)"`
	Name       string `orm:"size(255)"`
	Archived   bool
}

func init() {
	orm.RegisterModel(new(Category))
}

func (c *Category) TableName() string {
	return TYPE_TABLE_NAME
}

func (c *Category) Create() (err error) {
	if c.Name == "" {
		return fmt.Errorf("No name")
	}
	if c.LocationId <= 0 {
		return fmt.Errorf("No loc id")
	}

	o := orm.NewOrm()
	_, err = o.Insert(c)
	if err != nil {
		return fmt.Errorf("insert: %v", err)
	}

	return
}

func Get(id int64) (c *Category, err error) {
	o := orm.NewOrm()
	c = &Category{Id: id}
	err = o.Read(c)
	return
}

func GetAll(locId int64) (cs []*Category, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TYPE_TABLE_NAME).
		Filter("location_id", locId).
		All(&cs)
	return
}

func (c *Category) Archive() (err error) {
	o := orm.NewOrm()
	c.Archived = true
	_, err = o.Update(c)
	return
}

func (c *Category) Unarchive() (err error) {
	o := orm.NewOrm()
	c.Archived = false
	_, err = o.Update(c)
	return
}

func (c *Category) Update() (err error) {
	if c.LocationId <= 0 {
		return fmt.Errorf("No loc id")
	}
	o := orm.NewOrm()
	_, err = o.Update(c)
	return
}
