package invoices

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

const TABLE_NAME = "invoices"

func init() {
	orm.RegisterModel(new(Invoice))
}

// Invoice represents an actual or future invoice. Future invoices do not
// have a FastbillId.
type Invoice struct {
	Id          int64  `json:",omitempty"`
	LocationId  int64  `json:",omitempty"`
	FastbillId  int64  `json:",omitempty"`
	FastbillNo  string `json:",omitempty"`
	Month       int
	Year        int
	CustomerId  int64
	CustomerNo  int64
	UserId      int64
	Total       float64
	Status      string
	VatPercent  float64 `json:",omitempty"`
	Canceled    bool
	InvoiceDate time.Time
	PaidDate    time.Time
	DueDate     time.Time
}

func (inv *Invoice) TableName() string {
	return TABLE_NAME
}

// CreateOrUpdate based on LocationId and FastbillId
func CreateOrUpdate(invOrig *Invoice) (id int64, err error) {
	if invOrig.LocationId == 0 {
		return 0, fmt.Errorf("missing location id")
	}
	if invOrig.UserId == 0 {
		return 0, fmt.Errorf("missing user id")
	}
	if invOrig.Month == 0 {
		return 0, fmt.Errorf("missing month")
	}
	if invOrig.Year == 0 {
		return 0, fmt.Errorf("missing year")
	}
	o := orm.NewOrm()
	inv := *invOrig
	_, id, err = o.ReadOrCreate(&inv, "LocationId", "UserId", "Month", "Year")
	if err != nil {
		return 0, fmt.Errorf("read or create: %v", err)
	}

	inv = *invOrig
	inv.Id = id
	invOrig.Id = id
	if _, err = o.Update(&inv); err != nil {
		return inv.Id, fmt.Errorf("update: %v", err)
	}
	return
}

func CurrentInvoice(locationId, userId int64) (*Invoice, error) {
	inv := Invoice{
		LocationId: locationId,
		UserId:     userId,
		Month:      int(time.Now().Month()),
		Year:       time.Now().Year(),
		Status:     "draft",
	}
	_, _, err := orm.NewOrm().
		ReadOrCreate(&inv, "LocationId", "UserId", "Month", "Year", "Status")
	if err != nil {
		return nil, fmt.Errorf("read or create: %v", err)
	}
	return &inv, err
}

func Get(id int64) (*Invoice, error) {
	inv := Invoice{Id: id}
	err := orm.NewOrm().Read(&inv)
	return &inv, err
}

func GetAllInvoices(locId int64) ([]*Invoice, error) {
	var ivs []*Invoice
	_, err := orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locId).
		All(&ivs)
	return ivs, err
}

func GetAllInvoicesBetween(locId int64, year, month int) ([]*Invoice, error) {
	var ivs []*Invoice
	_, err := orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locId).
		Filter("year", year).
		Filter("month", month).
		All(&ivs)
	return ivs, err
}

func GetAllOfUserAt(locId, userId int64) ([]*Invoice, error) {
	var ivs []*Invoice
	_, err := orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locId).
		Filter("user_id", userId).
		All(&ivs)
	return ivs, err
}
