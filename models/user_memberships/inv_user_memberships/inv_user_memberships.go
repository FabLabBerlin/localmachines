package inv_user_memberships

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices"
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
	Id                    int64
	LocationId            int64
	UserId                int64
	MembershipId          int64
	UserMembershipId      int64
	UserMembership        *user_memberships.UserMembership `orm:"-" json:",omitempty"`
	StartDate             string
	TerminationDate       *string
	InitialDurationMonths int

	Created time.Time
	Updated time.Time

	InvoiceId     int64
	InvoiceStatus string
}

func Create(um *user_memberships.UserMembership, invoiceId int64) (
	ium *InvoiceUserMembership,
	err error,
) {
	ium = New(um, invoiceId)
	ium.Id, err = orm.NewOrm().Insert(ium)

	return
}

func New(um *user_memberships.UserMembership, invoiceId int64) *InvoiceUserMembership {
	return &InvoiceUserMembership{
		LocationId:            um.LocationId,
		UserId:                um.UserId,
		MembershipId:          um.MembershipId,
		UserMembershipId:      um.Id,
		UserMembership:        um,
		StartDate:             um.StartDate,
		TerminationDate:       um.TerminationDate,
		InitialDurationMonths: um.InitialDurationMonths,
		Created:               time.Now(),
		Updated:               time.Now(),
		InvoiceId:             invoiceId,
	}
}

var ErrNonDraftInvoice = errors.New("this would change a billed invoice")

func (this *InvoiceUserMembership) Denormalize(
	o orm.Ormer,
	um user_memberships.UserMembership,
) (err error) {
	inv, err := invoices.GetOrm(o, this.InvoiceId)
	if err != nil {
		return fmt.Errorf("get invoice: %v", err)
	}

	if inv.Status != "draft" {
		return ErrNonDraftInvoice
	}

	this.StartDate = um.StartDate
	this.TerminationDate = um.TerminationDate

	if err = this.Update(o); err != nil {
		return
	}

	return
}

func (this *InvoiceUserMembership) Membership() *memberships.Membership {
	return this.UserMembership.Membership
}

func (this *InvoiceUserMembership) TableName() string {
	return TABLE_NAME
}

func (this *InvoiceUserMembership) Update(o orm.Ormer) (err error) {
	_, err = o.Update(this)
	return
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

	ums, err := user_memberships.GetAllAt(locId)
	if err != nil {
		return
	}

	if deeplyPopulate(locId, iums, ums); err != nil {
		return nil, fmt.Errorf("deeply populate: %v", err)
	}

	return
}

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

	var ums []*user_memberships.UserMembership

	if _, err = orm.NewOrm().
		QueryTable("user_memberships").
		Filter("user_id", userId).
		All(&ums); err != nil {

		return
	}

	if deeplyPopulate(locId, iums, ums); err != nil {
		return nil, fmt.Errorf("deeply populate: %v", err)
	}

	return
}

func GetForUserMembership(userMembershipId int64) (iums []*InvoiceUserMembership, err error) {
	return GetForUserMembershipOrm(orm.NewOrm(), userMembershipId)
}

func GetForUserMembershipOrm(o orm.Ormer, userMembershipId int64) (iums []*InvoiceUserMembership, err error) {

	if _, err = o.
		QueryTable(TABLE_NAME).
		Filter("user_membership_id", userMembershipId).
		All(&iums); err != nil {

		return
	}

	if len(iums) == 0 {
		return
	}

	locId := iums[0].LocationId
	userId := iums[0].UserId

	var ums []*user_memberships.UserMembership

	if _, err = o.
		QueryTable("user_memberships").
		Filter("user_id", userId).
		All(&ums); err != nil {

		return
	}

	if deeplyPopulate(locId, iums, ums); err != nil {
		return nil, fmt.Errorf("deeply populate: %v", err)
	}

	return
}

func DeleteForUserMembership(o orm.Ormer, userMembershipId int64) (err error) {
	iums, err := GetForUserMembershipOrm(o, userMembershipId)
	if err != nil {
		return
	}

	if len(iums) == 0 {
		return
	}

	for _, ium := range iums {
		inv, err := invoices.GetOrm(o, ium.InvoiceId)
		if err != nil {
			return fmt.Errorf("get invoice: %v", err)
		}

		if inv.Status != "draft" {
			return ErrNonDraftInvoice
		}

		if _, err := o.Delete(ium); err != nil {
			return fmt.Errorf("delete: %v", err)
		}
	}

	return
}

func deeplyPopulate(
	locId int64,
	iums []*InvoiceUserMembership,
	ums []*user_memberships.UserMembership,
) (err error) {
	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		return fmt.Errorf("get all at: %v", err)
	}
	msbyId := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		msbyId[m.Id] = m
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

func (this *InvoiceUserMembership) StartDay() (d day.Day) {
	if d, err := day.NewString(this.StartDate); err == nil {
		return d
	}
	return
}

func (this *InvoiceUserMembership) TerminationDay() *day.Day {
	if td := this.TerminationDate; td != nil {
		d, err := day.NewString(*td)

		if err == nil {
			return &d
		} else {
			return nil
		}
	} else {
		return nil
	}
}
