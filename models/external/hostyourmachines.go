package external

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"regexp"
)

// Regular expression for email spec : RFC 5322
const _EXP_EMAIL = `(?i)[-a-z0-9~!$%^*_=+}{\'?]+(\.[-a-z0-9~!$%^*_=+}{\'?]+)*@([a-z0-9_][-a-z0-9_]*(\.[-a-z0-9_]+)*\.(aero|arpa|biz|com|coop|edu|gov|info|int|mil|museum|name|net|org|pro|travel|mobi|[a-z][a-z])|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,5})?`

func init() {
	orm.RegisterModel(new(Host))
}

type Host struct {
	Id           int64  `orm:"auto";"pk"`
	FirstName    string `orm:"size(100)"`
	LastName     string `orm:"size(100)"`
	Email        string `orm:"size(100)"`
	Location     string `orm:"size(100)"`
	Organization string `orm:"size(100)"`
	Phone        string `orm:"size(100)"`
	Comments     string `orm:"type(text)"`
}

func (this *Host) TableName() string {
	return "hosts"
}

func (this *Host) Add() (id int64, err error) {
	if this.FirstName == "" {
		return 0, fmt.Errorf("No first name")
	}

	if this.LastName == "" {
		return 0, fmt.Errorf("No last name")
	}

	if this.Email == "" {
		return 0, fmt.Errorf("No email")
	}

	exp, err := regexp.Compile(_EXP_EMAIL)
	if err != nil {
		return 0, fmt.Errorf("Failed to compile rexex: %v", err)
	}
	if !exp.MatchString(this.Email) {
		return 0, fmt.Errorf("Invalid email")
	}

	if this.Location == "" {
		return 0, fmt.Errorf("No location")
	}

	// TODO: Check for duplicates

	o := orm.NewOrm()
	id, err = o.Insert(this)
	return
}
