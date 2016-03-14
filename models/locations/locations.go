// locations of different Labs.
package locations

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/email"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

const TABLE_NAME = "locations"

type Location struct {
	Id               int64
	Title            string `orm:"size(100)"`
	FirstName        string `orm:"size(100)"`
	LastName         string `orm:"size(100)"`
	Email            string `orm:"size(100)"`
	City             string `orm:"size(100)"`
	Organization     string `orm:"size(100)"`
	Phone            string `orm:"size(100)"`
	Comments         string `orm:"type(text)"`
	Approved         bool
	XmppId           string `orm:"size(255)"`
	LocalIp          string `orm:"size(255)"`
	FeatureCoworking bool
	FeatureSetupTime bool
	FeatureSpaces    bool
	FeatureTutoring  bool
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
	if err != nil {
		return fmt.Errorf("insert: %v", err)
	}
	err = l.emailAnnounce()
	return
}

func (l *Location) emailAnnounce() (err error) {
	email := email.New()
	to := beego.AppConfig.String("mail@fablab.berlin")
	subject := "New location registered: " + l.Title
	message := "The location " + l.Title + " in " + l.City + " has been registered.\n\n"
	message += "The contact is: " + l.FirstName + " " + l.LastName + ",\n"
	message += l.Email + " (" + l.Phone + "). The location is not active and the approve flag\n"
	message += "must be set manually by the Database Administrator.\n"
	message += "Yours sincerely, Local Machines\n\n"
	return email.Send(to, subject, message)
}

func (l *Location) TableName() string {
	return TABLE_NAME
}

func Get(id int64) (l *Location, err error) {
	o := orm.NewOrm()
	l = &Location{Id: id}
	err = o.Read(l)
	return
}

func GetAll() (ls []*Location, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).All(&ls)
	return
}
