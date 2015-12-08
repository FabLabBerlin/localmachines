package purchases

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"time"
)

const (
	PURCHASE_TYPE_ACTIVATION     = "activation"
	PURCHASE_TYPE_CO_WORKING     = "co-working"
	PURCHASE_TYPE_RESERVATION    = "reservation"
	PURCHASE_TYPE_SPACE_PURCHASE = "space"
	PURCHASE_TYPE_TUTOR          = "tutor"
)

// This is a purchase row that appears in the XLSX file
type Purchase struct {
	Id        int64 `orm:"auto";"pk"`
	Type      string
	ProductId int64
	Created   time.Time `orm:"type(datetime)"`

	User   models.User `orm:"-" json:"-"`
	UserId int64

	TimeStart      time.Time `orm:"type(datetime)" json:",omitempty"`
	TimeEnd        time.Time `orm:"type(datetime)" json:",omitempty"`
	TimeEndPlanned time.Time `orm:"type(datetime)" json:",omitempty"`
	Quantity       float64
	PricePerUnit   float64
	PriceUnit      string
	Vat            float64
	Cancelled      bool

	TotalPrice      float64 `orm:"-"`
	DiscountedTotal float64 `orm:"-"`

	// Reservation fields:
	ReservationDisabled bool

	// Activation+Reservation fields:
	Machine   *models.Machine `orm:"-"`
	MachineId int64

	// Activation+Tutoring fields:
	Running bool

	// Old fields:
	Activation   *Activation          `orm:"-"`
	Reservation  *Reservation         `orm:"-"`
	MachineUsage time.Duration        `orm:"-"`
	Memberships  []*models.Membership `orm:"-"`

	Archived       bool
	Comments       string
	TimerTimeStart time.Time `orm:"type(datetime)" json:",omitempty"`
}

func (this *Purchase) TableName() string {
	return "purchases"
}

func init() {
	orm.RegisterModel(new(Purchase))
}

func GetPurchase(purchaseId int64) (purchase *Purchase, err error) {
	o := orm.NewOrm()
	purchase = &Purchase{Id: purchaseId}
	err = o.Read(purchase)
	return
}

func GetAllPurchasesOfType(purchaseType string) (purchases []*Purchase, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Purchase).TableName()).
		Filter("type", purchaseType).
		Exclude("archived", 1).
		All(&purchases)
	return
}

// Gets purchases of a specific user by consuming user ID.
func GetUserStartTime(userId int64) (startTime time.Time, err error) {
	query := "SELECT min(time_start) FROM purchases WHERE user_id = ?"
	o := orm.NewOrm()
	err = o.Raw(query, userId).QueryRow(&startTime)
	return
}

func DeletePurchase(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&Purchase{Id: id})
	return
}

func ArchivePurchase(purchase *Purchase) (err error) {
	o := orm.NewOrm()
	purchase.Archived = true
	_, err = o.Update(purchase)
	return
}

func (this *Purchase) MembershipStr() string {
	membershipStr := ""
	for _, membership := range this.Memberships {
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
		membershipStr = "None"
	}
	return membershipStr
}

func (this *Purchase) quantityFromTimes() (quantity float64) {
	var timeEnd time.Time
	if !this.TimeEnd.IsZero() {
		timeEnd = this.TimeEnd
	} else if !this.TimeEndPlanned.IsZero() {
		timeEnd = this.TimeEndPlanned
	} else {
		timeEnd = time.Now()
	}

	seconds := timeEnd.Sub(this.TimeStart).Seconds()

	switch this.PriceUnit {
	case "minute":
		quantity = float64(seconds) / 60
	case "30 minutes":
		quantity = float64(seconds) / 1800
	case "hour":
		quantity = float64(seconds) / 3600
	default:
		beego.Error("unknown price unit ", this.PriceUnit)
	}

	return
}

func PriceTotalExclDisc(p *Purchase) float64 {
	return p.Quantity * p.PricePerUnit
}

func PriceTotalDisc(p *Purchase) (float64, error) {
	if p.Type != PURCHASE_TYPE_ACTIVATION {
		return PriceTotalExclDisc(p), nil
	}

	priceTotal := PriceTotalExclDisc(p)
	for _, membership := range p.Memberships {

		// We need to know whether the machine is affected by the base membership
		// as well as the individual activation is affected by the user membership
		isAffected, err := membership.IsMachineAffected(p.Machine.Id)
		if err != nil {
			beego.Error(
				"Failed to check whether machine is affected by membership:", err)
			return 0, fmt.Errorf(
				"Failed to check whether machine is affected by membership")
		}

		machinePriceDeduction := 0.0
		if isAffected {
			machinePriceDeduction = float64(membership.MachinePriceDeduction)
		}
		// Discount total price
		priceTotal = priceTotal - (priceTotal * machinePriceDeduction / 100.0)
	}
	return priceTotal, nil
}

type Purchases struct {
	Data []*Purchase
}

func (this Purchases) Len() int {
	return len(this.Data)
}

func (this Purchases) Less(i, j int) bool {
	var timeStartI = (*this.Data[i]).TimeStart
	var timeStartJ = (*this.Data[j]).TimeStart
	if timeStartI.Before(timeStartJ) {
		return true
	} else if timeStartJ.Before(timeStartI) {
		return false
	} else if (*this.Data[i]).Machine != nil && (*this.Data[j]).Machine != nil {
		return (*this.Data[i]).Machine.Name < (*this.Data[j]).Machine.Name
	} else {
		return false
	}
}

func (this Purchase) ProductName() string {
	switch this.Type {
	case PURCHASE_TYPE_ACTIVATION:
		return this.Machine.Name
	case PURCHASE_TYPE_RESERVATION:
		return "Reservation (" + this.Machine.Name + ")"
	}
	return "Unnamed product"
}

func (this Purchases) Swap(i, j int) {
	*this.Data[i], *this.Data[j] = *this.Data[j], *this.Data[i]
}
