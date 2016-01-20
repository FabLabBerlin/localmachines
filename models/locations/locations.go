package locations

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
)

const TABLE_NAME = "locations"

type Location struct {
	Id           int64
	Title        string `orm:"size(100)"`
	FirstName    string `orm:"size(100)"`
	LastName     string `orm:"size(100)"`
	Email        string `orm:"size(100)"`
	City         string `orm:"size(100)"`
	Organization string `orm:"size(100)"`
	Phone        string `orm:"size(100)"`
	Comments     string `orm:"type(text)"`
	Approved     bool
}

func init() {
	orm.RegisterModel(new(Location))
}

func (l *Location) Save() (err error) {
	if l.Title == "" {
		return fmt.Errorf("No title")
	}

	if l.FirstName == "" {
		return fmt.Errorf("No first name")
	}

	if l.LastName == "" {
		return fmt.Errorf("No last name")
	}

	if l.Email == "" {
		return fmt.Errorf("No email")
	}

	// Just a rough check for email because addresses like
	// foo+bar@fablab.berlin are difficult to test
	if !strings.Contains(l.Email, "@") || len(l.Email) < 5 {
		return fmt.Errorf("Invalid email")
	}

	if l.City == "" {
		return fmt.Errorf("No city")
	}

	// TODO: Check for duplicates

	o := orm.NewOrm()
	_, err = o.Insert(l)
	return
}

func (l *Location) TableName() string {
	return TABLE_NAME
}

func GetAll() (ls []*Location, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).All(&ls)
	return
}
