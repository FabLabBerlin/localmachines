package invoices

import (
	"fmt"
	"github.com/astaxie/beego"
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
	beego.Info("CreateOrUpdate: invOrig.UserId=", invOrig.UserId)
	if invOrig.LocationId == 0 {
		return 0, fmt.Errorf("missing location id")
	}
	if invOrig.FastbillId == 0 {
		return 0, fmt.Errorf("missing fastbill id")
	}
	o := orm.NewOrm()
	inv := *invOrig
	created, id, err := o.ReadOrCreate(&inv, "LocationId", "FastbillId")
	if err != nil {
		return id, fmt.Errorf("read or create: %v", err)
	}
	if inv.UserId == 19 {
		beego.Info("CreateOrUpdate: created=", created)
		beego.Info("CreateOrUpdate: id=", id)
		beego.Info("CreateOrUpdate: inv=", inv)
	}
	inv = *invOrig
	inv.Id = id
	if _, err = o.Update(&inv); err != nil {
		return inv.Id, fmt.Errorf("update: %v", err)
	}
	return
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
