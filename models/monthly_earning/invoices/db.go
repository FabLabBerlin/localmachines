package invoices

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Invoice))
}

func (inv *Invoice) TableName() string {
	return TABLE_NAME
}

// CreateOrUpdate based on LocationId and FastbillId
func CreateOrUpdate(invOrig *Invoice) (id int64, err error) {
	if invOrig.LocationId == 0 {
		return 0, fmt.Errorf("missing location id")
	}
	if invOrig.FastbillId == 0 {
		return 0, fmt.Errorf("missing fastbill id")
	}
	o := orm.NewOrm()
	inv := *invOrig
	_, id, err = o.ReadOrCreate(&inv, "LocationId", "FastbillId")
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
