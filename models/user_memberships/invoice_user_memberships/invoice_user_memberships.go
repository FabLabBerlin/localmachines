package invoice_user_memberships

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type InvoiceUserMembership struct {
	Id               int64
	UserId           int64
	MembershipId     int64
	UserMembershipId int64
	StartDate        time.Time `orm:"type(datetime)"`
	TerminationDate  time.Time `orm:"type(datetime)"`
	DurationMonths   int64

	Created time.Time
	Updated time.Time

	InvoiceId     int64
	InvoiceStatus string
}

func (this *InvoiceUserMembership) TableName() string {
	return "invoice_user_memberships"
}

func init() {
	orm.RegisterModel(new(InvoiceUserMembership))
}
