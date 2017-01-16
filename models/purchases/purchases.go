package purchases

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

const TABLE_NAME = "purchases"

const (
	TYPE_ACTIVATION  = "activation"
	TYPE_OTHER       = "other"
	TYPE_RESERVATION = "reservation"
	TYPE_TUTOR       = "tutor"
)

// This is a purchase row that appears in the XLSX file
type Purchase struct {
	Id         int64
	LocationId int64
	Type       string
	ProductId  int64
	Created    time.Time `orm:"type(datetime)"`

	User   users.User `orm:"-" json:"-"`
	UserId int64

	TimeStart      time.Time `orm:"type(datetime)" json:",omitempty"`
	TimeEndPlanned time.Time `orm:"type(datetime)" json:",omitempty"`
	Quantity       float64
	PricePerUnit   float64
	PriceUnit      string
	Vat            float64
	Cancelled      bool

	TotalPrice      float64 `orm:"-"`
	DiscountedTotal float64 `orm:"-"`
	PriceVAT        float64 `orm:"-"`
	PriceExclVAT    float64 `orm:"-"`

	InvoiceId     int64
	InvoiceStatus string

	// Reservation fields:
	ReservationDisabled bool

	// Activation+Reservation fields:
	Machine   *machine.Machine `orm:"-"`
	MachineId int64

	// Type Other fields:
	CustomName string

	// Activation+Tutoring fields:
	Running bool

	// Old fields:
	Memberships []*memberships.Membership `orm:"-"`

	Archived       bool
	Comments       string
	TimerTimeStart time.Time `orm:"type(datetime)" json:",omitempty"`
}

func (this *Purchase) TableName() string {
	return TABLE_NAME
}

func init() {
	orm.RegisterModel(new(Purchase))
}

func Create(p *Purchase) (err error) {
	p.Created = time.Now()
	p.Id, err = CreateOrm(orm.NewOrm(), p)
	return
}

func CreateOrm(o orm.Ormer, p *Purchase) (id int64, err error) {
	if p.LocationId <= 0 {
		return 0, errors.New("LocationId must be > 0")
	}
	if p.InvoiceId <= 0 {
		return 0, errors.New("InvoiceId must be > 0")
	}
	if id, err = o.Insert(p); err != nil {
		return
	}
	p.Id = id
	return
}

func Get(purchaseId int64) (purchase *Purchase, err error) {
	o := orm.NewOrm()
	purchase = &Purchase{Id: purchaseId}
	err = o.Read(purchase)
	return
}

func GetAllAt(locationId int64) (ps []*Purchase, err error) {
	o := orm.NewOrm()

	_, err = o.QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		Exclude("invoice_status", "outgoing").
		Exclude("invoice_status", "credit").
		Exclude("archived", "1").
		All(&ps)

	return

}

func GetByInvoiceId(invoiceId int64) (ps []*Purchase, err error) {
	if invoiceId <= 0 {
		return nil, fmt.Errorf("invoice id is %v", invoiceId)
	}

	o := orm.NewOrm()

	_, err = o.QueryTable(TABLE_NAME).
		Filter("invoice_id", invoiceId).
		Exclude("archived", "1").
		All(&ps)

	return
}

func GetAllBetweenAt(locationId int64, interval lib.Interval) (ps []*Purchase, err error) {
	all, err := GetAllAt(locationId)
	if err != nil {
		return nil, fmt.Errorf("get all at: %v", err)
	}

	ps = make([]*Purchase, 0, len(all))
	for _, p := range all {
		if interval.Contains(p.TimeStart) && !p.Running {
			if p.UserId != 0 {
				ps = append(ps, p)
			} else {
				beego.Warn("UserId = 0 for Purchase # ", p.Id)
			}
		}
	}

	return
}

func GetAllOfType(purchaseType string) (purchases []*Purchase, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("type", purchaseType).
		Exclude("invoice_status", "outgoing").
		Exclude("invoice_status", "credit").
		Exclude("archived", 1).
		All(&purchases)
	return
}

func GetAllOfTypeAt(locationId int64, typ string) (ps []*Purchase, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		Filter("type", typ).
		Exclude("invoice_status", "outgoing").
		Exclude("invoice_status", "credit").
		Exclude("archived", 1).
		All(&ps)
	return
}

// Gets purchases of a specific user by consuming user ID.
func GetUserStartTime(userId int64) (startTime time.Time, err error) {
	query := "SELECT min(time_start) FROM purchases WHERE user_id = ?"
	o := orm.NewOrm()
	err = o.Raw(query, userId).QueryRow(&startTime)
	return
}

func Delete(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&Purchase{Id: id})
	return
}

func Archive(purchase *Purchase) (err error) {
	o := orm.NewOrm()
	purchase.Archived = true
	_, err = o.Update(purchase)
	return
}

func (this Purchase) CloneOrm(o orm.Ormer, inv *invoices.Invoice) error {
	var clone Purchase
	clone = this
	clone.Id = 0
	clone.InvoiceId = inv.Id
	clone.InvoiceStatus = inv.Status
	_, err := CreateOrm(o, &clone)
	return err
}

func (this *Purchase) Duration() time.Duration {
	return time.Duration(this.Seconds()) * time.Second
}

func (this *Purchase) MembershipStr() (membershipStr string, err error) {
	affected, err := AffectedMemberships(this)
	if err != nil {
		return "", fmt.Errorf("affected memberships: %v", err)
	}
	for _, membership := range affected {
		memStr := fmt.Sprintf("%s (%d%%)",
			membership.ShortName,
			membership.MachinePriceDeduction)
		if membershipStr == "" {
			membershipStr = memStr
		} else {
			membershipStr = fmt.Sprintf("%s, %s",
				membershipStr, memStr)
		}
	}
	if membershipStr == "" {
		membershipStr = "Pay-As-You-Go"
	}
	return
}

func (this *Purchase) quantityFromTimes(timeEnd time.Time) (quantity float64) {
	seconds := timeEnd.Sub(this.TimeStart).Seconds()

	switch this.PriceUnit {
	case "minute":
		return float64(seconds) / 60
	case "30 minutes":
		return float64(seconds) / 1800
	case "hour":
		return float64(seconds) / 3600
	default:
		beego.Error("unknown price unit ", this.PriceUnit, " for #", this.Id)
	}

	return
}

func (this *Purchase) Seconds() (seconds float64) {
	if this.Quantity == 0 {
		return 0
	} else {
		switch this.PriceUnit {
		case "second":
			return this.Quantity
		case "minute":
			return float64(this.Quantity) * 60
		case "30 minutes":
			return float64(this.Quantity) * 60 * 30
		case "hour":
			return float64(this.Quantity) * 60 * 60
		case "day":
			return float64(this.Quantity) * 60 * 60 * 24
		case "month":
			return float64(this.Quantity) * 60 * 60 * 24 * 30
		default:
			beego.Error("unknown price unit ", this.PriceUnit)
		}
	}
	return
}

func Update(p *Purchase) (err error) {
	o := orm.NewOrm()

	inv, err := invoices.Get(p.InvoiceId)
	if err != nil {
		return fmt.Errorf("get invoice: %v", err)
	}

	if inv.Status != "draft" {
		return fmt.Errorf("invoice doesn't have status draft: %v", inv.Status)
	}
	p.InvoiceStatus = inv.Status

	_, err = o.Update(p)
	return
}

func PriceTotalExclDisc(p *Purchase) float64 {
	return p.Quantity * p.PricePerUnit
}

func PriceTotalDisc(p *Purchase) (float64, error) {
	if p.Type != TYPE_ACTIVATION {
		return PriceTotalExclDisc(p), nil
	}

	priceTotal := PriceTotalExclDisc(p)
	affectedMemberships, err := AffectedMemberships(p)
	if err != nil {
		return 0, fmt.Errorf("affected memberships: %v", err)
	}
	fmt.Printf("PriceTotalDisc: affectedMemberships=%v\n", affectedMemberships)
	for _, membership := range affectedMemberships {
		machinePriceDeduction := float64(membership.MachinePriceDeduction)
		// Discount total price
		priceTotal = priceTotal - (priceTotal * machinePriceDeduction / 100.0)
	}
	return priceTotal, nil
}

func AffectedMemberships(p *Purchase) (affected []*memberships.Membership, err error) {
	fmt.Printf("purchase.go:AffectedMemberships()\n")
	affected = make([]*memberships.Membership, 0, 2)

	if p.MachineId == 0 {
		return
	}

	for _, membership := range p.Memberships {

		isAffected := false /*, err := membership.IsMachineAffected(p.MachineId)
		if err != nil {
			return nil, fmt.Errorf(
				"check if machine affected by membership %v: %v", membership.Id, err)
		}*/

		categoryIds, err := membership.AffectedCategoryIds()
		if err != nil {
			return nil, fmt.Errorf("affected category ids: %v", err)
		}

		for _, id := range categoryIds {
			if p.Machine.TypeId == id {
				isAffected = true
			}
		}

		fmt.Printf("AffectedMemberships: isAffected=%v\n", isAffected)

		if isAffected {
			affected = append(affected, membership)
		}
	}

	return
}

type Purchases []*Purchase

func (this Purchases) Len() int {
	return len(this)
}

func (this Purchases) Less(i, j int) bool {
	var timeStartI = (*this[i]).TimeStart
	var timeStartJ = (*this[j]).TimeStart
	if timeStartI.Before(timeStartJ) {
		return true
	} else if timeStartJ.Before(timeStartI) {
		return false
	} else if (*this[i]).Machine != nil && (*this[j]).Machine != nil {
		return (*this[i]).Machine.Name < (*this[j]).Machine.Name
	} else {
		return false
	}
}

func (this Purchase) ProductName() string {
	switch this.Type {
	case TYPE_ACTIVATION:
		if this.Machine != nil {
			return this.Machine.Name
		}
	case TYPE_RESERVATION:
		if this.Machine != nil {
			return this.Machine.Name + " Reservation"
		}
	case TYPE_TUTOR:
		return "Tutoring"
	case TYPE_OTHER:
		return this.CustomName
	}

	return "Unnamed product"
}

func (this Purchases) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

func (this Purchase) TimeEnd() time.Time {
	return this.TimeStart.Add(time.Duration(this.Seconds()) * time.Second)
}
