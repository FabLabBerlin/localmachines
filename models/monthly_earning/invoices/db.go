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
func CreateOrUpdate(inv *Invoice) (id int64, err error) {
	if inv.LocationId == 0 {
		return 0, fmt.Errorf("missing location id")
	}
	if inv.FastbillId == 0 {
		return 0, fmt.Errorf("missing fastbill id")
	}
	o := orm.NewOrm()
	_, id, err = o.ReadOrCreate(inv, "LocationId", "FastbillId")
	if err != nil {
		return id, fmt.Errorf("read or create: %v", err)
	}
	inv.Id = id
	if _, err = o.Update(inv); err != nil {
		return inv.Id, fmt.Errorf("update: %v", err)
	}
	return
}
