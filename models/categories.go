package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Category))
}

type Category struct {
	Id                int64  `orm:"auto";"pk"`
	Name              string `orm:"size(255)"`
	Shortname         string `orm:"size(100)"`
	Description       string `orm:"type(text)"`
	Image             string `orm:"size(255)"` // TODO: media and media type tables
}

func (c *Category) TableName() string {
	return "categories"
}

func GetAllCategories() (categories []Category, err error) {
	o := orm.NewOrm()
	c := Category{}
	_, err = o.QueryTable(c.TableName()).All(&categories)
	return
}

func CreateCategory(name string) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(&Category{Name: name})
}

func UpdateCategory(category *Category) error {
	o := orm.NewOrm()
	_,err := o.Update(category)
	return err
}

func DeleteCategory(categoryId int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Category{Id: categoryId})
	return err
}

	

