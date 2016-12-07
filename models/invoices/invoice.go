/*
invoices primitives.

Building blocks for all invoicing functionality, mostly based on the
Invoice type.
*/
package invoices

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
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
	Id                 int64  `json:",omitempty"`
	LocationId         int64  `json:",omitempty"`
	FastbillId         int64  `json:",omitempty"`
	FastbillNo         string `json:",omitempty"`
	CanceledFastbillId int64  `json:",omitempty"`
	CanceledFastbillNo string `json:",omitempty"`
	Month              int
	Year               int
	CustomerId         int64
	CustomerNo         int64
	UserId             int64
	Total              float64
	Status             string
	VatPercent         float64 `json:",omitempty"`
	Canceled           bool
	Sent               bool
	CanceledSent       bool
	InvoiceDate        time.Time
	PaidDate           time.Time
	DueDate            time.Time
	Current            bool
}

func Create(inv *Invoice) (id int64, err error) {
	if err := inv.assertDataOk(); err != nil {
		return 0, fmt.Errorf("assert data ok: %v", err)
	}

	o := orm.NewOrm()

	if err = o.Begin(); err != nil {
		return 0, fmt.Errorf("begin tx: %v", err)
	}

	if id, err = o.Insert(inv); err != nil {
		return 0, fmt.Errorf("insert inv: %v", err)
	}

	if now := time.Now(); inv.Month == int(now.Month()) &&
		inv.Year == now.Year() {
		if inv.setCurrent(o); err != nil {
			o.Rollback()
			return 0, fmt.Errorf("set current: %v", err)
		}
	}

	if err = o.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %v", err)
	}

	return
}

// ErrNoInvoiceForThatMonth is related to lazy creation of invoices. We don't
// want to spam the database with irrelevant data.
var ErrNoInvoiceForThatMonth = errors.New("no draft invoice exists for that month (and none will be auto-created)")

// GetDraft for User uid @ Location locId and time t.  If it doesn't exist, it
// gets created insofar it doesn't violate business logic.
func GetDraft(locId, uid int64, t time.Time) (*Invoice, error) {
	y := t.Year()
	m := t.Month()
	isThisMonth := time.Now().Month() == m && time.Now().Year() == y

	timeNextMonth := time.Now().AddDate(0, 1, 0)
	isNextMonth := timeNextMonth.Month() == m && timeNextMonth.Year() == y

	inv := Invoice{
		LocationId: locId,
		UserId:     uid,
		Month:      int(m),
		Year:       y,
		Status:     "draft",
	}

	existing, err := getByProps(
		inv.LocationId,
		inv.UserId,
		time.Month(inv.Month),
		inv.Year,
		inv.Status,
	)

	if err == nil {
		return existing, nil
	} else if err == orm.ErrNoRows {
		if isThisMonth || isNextMonth {
			if _, err := Create(&inv); err == nil {
				if isThisMonth {
					if err := inv.SetCurrent(); err != nil {
						return nil, fmt.Errorf("set current: %v", err)
					}
				}
				return &inv, nil
			} else {
				return nil, fmt.Errorf("create: %v", err)
			}
		} else {
			return nil, ErrNoInvoiceForThatMonth
		}
	} else {
		return nil, fmt.Errorf("get by props: %v", err)
	}
}

func Get(id int64) (*Invoice, error) {
	return GetOrm(orm.NewOrm(), id)
}

func GetOrm(o orm.Ormer, id int64) (*Invoice, error) {
	inv := Invoice{Id: id}
	err := o.Read(&inv)
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

func getByProps(locId, uid int64, m time.Month, y int, status string) (*Invoice, error) {
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
			iv.Status == inv.Status &&
			iv.Canceled == inv.Canceled &&
			iv.CanceledFastbillNo == inv.CanceledFastbillNo {
			return fmt.Errorf("(id=%v) conflicting with invoice %v",
				inv.Id, iv.Id)
		}
	}

	return
}

func (inv *Invoice) Interval() lib.Interval {
	if inv.Month == 0 || inv.Year == 0 {
		panic(fmt.Sprintf("inv.Month=%v, inv.Year=%v", inv.Month, inv.Year))
	}

	return lib.Interval{
		MonthFrom: inv.Month,
		YearFrom:  inv.Year,
		MonthTo:   inv.Month,
		YearTo:    inv.Year,
	}
}

func (inv *Invoice) Save() (err error) {
	return inv.SaveOrm(orm.NewOrm())
}

func (inv *Invoice) SaveOrm(o orm.Ormer) (err error) {
	if err := inv.assertDataOk(); err != nil {
		return fmt.Errorf("assert data ok: %v", err)
	}
	_, err = o.Update(inv)
	return
}

func (inv *Invoice) SaveTotal() (err error) {
	if err := inv.assertDataOk(); err != nil {
		return fmt.Errorf("assert data ok: %v", err)
	}
	if inv.Id <= 0 {
		return fmt.Errorf("no (valid) invoice id")
	}

	query := `
	UPDATE invoices
	SET total = ?
	WHERE id = ?
	`

	o := orm.NewOrm()
	_, err = o.Raw(query, inv.Total, inv.Id).Exec()

	return
}

// SetCurrent sets inv as the current invoice.  The function is idempotent.
// To avoid expensive transactions, it does these steps:
//
// 1. If there is already another current invoice for user @ location:
//
// 1a. Get all memberships from it
//
// 1b. Clone them for the new invoice, if all aren't already present
//
// 2. Transactionally switch over current=true state
func (inv *Invoice) SetCurrent() (err error) {
	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		return fmt.Errorf("begin tx: %v", err)
	}

	if err := inv.setCurrent(o); err != nil {
		o.Rollback()
		return err
	}

	if err := o.Commit(); err != nil {
		return fmt.Errorf("commit tx: %v", err)
	}

	return
}

func (inv *Invoice) setCurrent(o orm.Ormer) (err error) {
	if inv.Current {
		return nil
	}

	var currentInvoice *Invoice

	var tmp Invoice
	err = o.QueryTable(TABLE_NAME).
		Filter("location_id", inv.LocationId).
		Filter("user_id", inv.UserId).
		Filter("current", true).
		One(&tmp)

	if err == nil {
		currentInvoice = &tmp
	} else if err == orm.ErrNoRows {
		err = nil
	} else {
		return fmt.Errorf("one: %v", err)
	}

	if currentInvoice != nil {
		currentInvoice.Current = false
		if _, err := o.Update(currentInvoice); err != nil {
			return fmt.Errorf("update current invoice: %v", err)
		}
	}

	inv.Current = true
	if _, err := o.Update(inv); err != nil {
		return fmt.Errorf("update this invoice: %v", err)
	}

	return
}

func (inv *Invoice) TableName() string {
	return TABLE_NAME
}
