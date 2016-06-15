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
	Current     bool
}

func (inv *Invoice) TableName() string {
	return TABLE_NAME
}

func Create(inv *Invoice) (id int64, err error) {
	if err := inv.assertDataOk(); err != nil {
		return 0, fmt.Errorf("assert data ok: %v", err)
	}

	o := orm.NewOrm()
	_, err = o.Insert(inv)

	return
}

// CreateOrUpdate based on LocationId and FastbillId
func CreateOrUpdate(invOrig *Invoice) (id int64, err error) {
	if err := invOrig.assertDataOk(); err != nil {
		return 0, fmt.Errorf("assert data ok: %v", err)
	}

	existing, err := GetByProps(
		invOrig.LocationId,
		invOrig.UserId,
		time.Month(invOrig.Month),
		invOrig.Year,
		invOrig.Status,
	)
	if err == nil {
		invOrig.Id = existing.Id
		if err = invOrig.Update(); err == nil {
			return existing.Id, nil
		} else {
			return 0, fmt.Errorf("update: %v", err)
		}
	} else if err == orm.ErrNoRows {
		id, err := Create(invOrig)
		if err == nil {
			return id, err
		} else {
			return 0, fmt.Errorf("create: %v", err)
		}
	} else {
		return 0, fmt.Errorf("get by props: %v", err)
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

	existing, err := GetByProps(
		inv.LocationId,
		inv.UserId,
		time.Month(inv.Month),
		inv.Year,
		inv.Status,
	)

	if err == nil {
		return existing, nil
	} else if err == orm.ErrNoRows {
		if _, err := Create(&inv); err == nil {
			return &inv, nil
		} else {
			return nil, fmt.Errorf("create: %v", err)
		}
	} else {
		return nil, fmt.Errorf("get by props: %v", err)
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

func GetByProps(locId, uid int64, m time.Month, y int, status string) (*Invoice, error) {
	var inv Invoice

	err := orm.NewOrm().
		QueryTable(TABLE_NAME).
		Filter("location_id", locId).
		Filter("user_id", uid).
		Filter("month", int(m)).
		Filter("year", y).
		Filter("status", status).
		One(&inv)

	return &inv, err
}

func (inv *Invoice) assertDataOk() (err error) {
	// Basic checks
	if inv.LocationId == 0 {
		return fmt.Errorf("missing location id")
	}
	if inv.UserId == 0 {
		return fmt.Errorf("missing user id")
	}
	if inv.Month == 0 {
		return fmt.Errorf("missing month")
	}
	if inv.Year == 0 {
		return fmt.Errorf("missing year")
	}
	if inv.Status == "" {
		return fmt.Errorf("missing status")
	}

	// Check for possible conflicts
	ivs, err := GetAllOfUserAt(inv.LocationId, inv.UserId)
	if err != nil {
		return fmt.Errorf("get all of user at: %v", err)
	}
	for _, iv := range ivs {
		if iv.Id == inv.Id {
			continue
		}

		if iv.LocationId == inv.LocationId &&
			iv.UserId == inv.UserId &&
			iv.Month == inv.Month &&
			iv.Year == inv.Year &&
			iv.Status == inv.Status {
			return fmt.Errorf("(id=%v) conflicting with invoice %v",
				inv.Id, iv.Id)
		}
	}

	return
}

func (inv *Invoice) Update() (err error) {
	if err := inv.assertDataOk(); err != nil {
		return fmt.Errorf("assert data ok: %v", err)
	}
	_, err = orm.NewOrm().Update(inv)
	return
}
