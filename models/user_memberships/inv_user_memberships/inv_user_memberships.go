package inv_user_memberships

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/astaxie/beego/orm"
	"time"
)

const TABLE_NAME = "invoice_user_memberships"

// InvoiceUserMembership represents the billed user memberships. As such it is
// not directly user editable. The only possibility is editing UserMembership,
// in case this entity is associated with a draft invoices, the changes will be
// propagated.
type InvoiceUserMembership struct {
	Id               int64
	LocationId       int64
	UserId           int64
	MembershipId     int64
	UserMembershipId int64
	UserMembership   *user_memberships.UserMembership `orm:"-" json:",omitempty"`
	StartDate        time.Time                        `orm:"type(datetime)"`
	TerminationDate  time.Time                        `orm:"type(datetime)"`
	DurationMonths   int64

	Created time.Time
	Updated time.Time

	InvoiceId     int64
	InvoiceStatus string
}

func (this *InvoiceUserMembership) Membership() *memberships.Membership {
	return this.UserMembership.Membership
}

func (this *InvoiceUserMembership) TableName() string {
	return TABLE_NAME
}

func init() {
	orm.RegisterModel(new(InvoiceUserMembership))
}

func GetAllAt(locId int64) (iums []*InvoiceUserMembership, err error) {

	if _, err = orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locId).
		All(&iums); err != nil {

		return
	}

	if len(iums) == 0 {
		return
	}

	ums, err := user_memberships.GetAllAtDeep(locId)
	if err != nil {
		return
	}
	umsById := make(map[int64]*user_memberships.UserMembership)
	for _, um := range ums {
		umsById[um.Id] = um
	}

	for _, ium := range iums {
		ium.UserMembership = umsById[ium.UserMembershipId]
	}

	return
}

// Gets all user memberships for a user by consuming user ID.
func GetForInvoice(invoiceId int64) (iums []*InvoiceUserMembership, err error) {

	if _, err = orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("invoice_id", invoiceId).
		All(&iums); err != nil {

		return
	}

	if len(iums) == 0 {
		return
	}

	locId := iums[0].LocationId
	userId := iums[0].UserId

	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all at: %v", err)
	}
	msbyId := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		msbyId[m.Id] = m
	}

	var ums []*user_memberships.UserMembership

	if _, err = orm.NewOrm().
		QueryTable("user_memberships").
		Filter("user_id", userId).
		All(&ums); err != nil {

		return
	}
	umsById := make(map[int64]*user_memberships.UserMembership)
	for _, um := range ums {
		um.Membership = msbyId[um.MembershipId]
		umsById[um.Id] = um
	}

	for _, ium := range iums {
		ium.UserMembership = umsById[ium.UserMembershipId]
	}

	return
}
