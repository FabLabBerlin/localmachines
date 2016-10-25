package invoice_user_memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/astaxie/beego/orm"
	"time"
)

type InvoiceUserMembership struct {
	Id               int64
	UserId           int64
	MembershipId     int64
	UserMembershipId int64
	StartDate        time.Time `orm:"type(datetime)"`
	EndDate          time.Time `orm:"type(datetime)"`

	InvoiceId     int64
	InvoiceStatus string
}

func (this *InvoiceUserMembership) TableName() string {
	return "invoice_user_memberships"
}

func init() {
	orm.RegisterModel(new(InvoiceUserMembership))
}
