package invoices

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/memberships"
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
	Sent        bool
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

	if err = o.Begin(); err != nil {
		return 0, fmt.Errorf("begin tx: %v", err)
	}

	if _, err = o.Insert(inv); err != nil {
		return 0, fmt.Errorf("insert inv: %v", err)
	}

	if now := time.Now(); inv.Month == int(now.Month()) &&
		inv.Year == now.Year() {
		lastMonth := now.AddDate(0, -1, 0)
		lastInv, err := InvoiceOfMonth(inv.LocationId, inv.UserId, lastMonth.Year(), lastMonth.Month())
		if err == nil {
			umbs, err := memberships.GetUserMembershipsForInvoice(lastInv.Id)
			if err != nil {
				return 0, fmt.Errorf("get user memberships for last month: %v", err)
			}
			for _, umb := range umbs.Data {
				if !umb.AutoExtend &&
					umb.EndDate.Before(inv.Interval().TimeFrom()) {
					continue
				}

				clone := umb.UserMembership()
				clone.Id = 0
				clone.InvoiceId = inv.Id
				clone.InvoiceStatus = inv.Status

				if _, err := o.Insert(&clone); err != nil {
					return 0, fmt.Errorf("inserting cloned membership: %v", err)
				}
			}
		} else if err == ErrNoInvoiceForThatMonth {
			// Nothing to do:
			err = nil
		} else {
			return 0, fmt.Errorf("getting in inv of last month: %v", err)
		}
	}

	if err = o.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %v", err)
	}

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
		if err = invOrig.Save(); err == nil {
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
}

func ThisMonthInvoice(locId, uid int64) (*Invoice, error) {
	return InvoiceOfMonth(locId, uid, time.Now().Year(), time.Now().Month())
}

var ErrNoInvoiceForThatMonth = errors.New("no invoice exists for that month (and none will be auto-created)")

func InvoiceOfMonth(locId, uid int64, y int, m time.Month) (*Invoice, error) {
	isThisMonth := time.Now().Month() == m && time.Now().Year() == y

	inv := Invoice{
		LocationId: locId,
		UserId:     uid,
		Month:      int(m),
		Year:       y,
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
		if isThisMonth {
			if _, err := Create(&inv); err == nil {
				if err := inv.SetCurrent(); err != nil {
					return nil, fmt.Errorf("set current: %v", err)
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

func (inv *Invoice) AttachUserMembership(um *memberships.UserMembership) error {
	if inv.Id == 0 {
		return errors.New("invoice Id = 0")
	}

	switch um.InvoiceId {
	case inv.Id:
		return nil
	case 0:
		um.InvoiceId = inv.Id
		if err := um.Update(); err != nil {
			return fmt.Errorf("update user membership: %v", err)
		}
		return nil
	default:
		locId := inv.LocationId
		ums, err := memberships.GetAllUserMembershipsAt(locId)
		if err != nil {
			return fmt.Errorf("get all user memberships at %v: %v", locId, err)
		}
		for _, existing := range ums {
			if existing.InvoiceId == inv.Id {
				// Already done
				return nil
			}
		}

		o := orm.NewOrm()
		newUm, err := memberships.CreateUserMembership(o, um.UserId, um.MembershipId, inv.Id, um.StartDate)
		if err != nil {
			return fmt.Errorf("create user membership: %v", err)
		}
		newUm.InvoiceStatus = inv.Status
		if newUm.Update(); err != nil {
			return fmt.Errorf("update user membership: %v", err)
		}
	}
	return nil
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
	if err := inv.assertDataOk(); err != nil {
		return fmt.Errorf("assert data ok: %v", err)
	}
	_, err = orm.NewOrm().Update(inv)
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
// 1a. Get all memberships from it
// 1b. Clone them for the new invoice, if all aren't already present
// 2. Transactionally switch over current=true state
func (inv *Invoice) SetCurrent() (err error) {
	if inv.Current {
		return nil
	}

	var currentInvoice *Invoice

	o := orm.NewOrm()

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
		currentUms, err := memberships.GetUserMembershipsForInvoice(currentInvoice.Id)
		if err != nil {
			return fmt.Errorf("get user memberships for current inv: %v", err)
		}

		thisUms, err := memberships.GetUserMembershipsForInvoice(inv.Id)
		if err != nil {
			return fmt.Errorf("get user memberships for this inv: %v", err)
		}

		umsToBeCloned := make([]*memberships.UserMembershipCombo, 0, len(currentUms.Data))
		for _, um := range currentUms.Data {
			alreadyCloned := false
			for _, existing := range thisUms.Data {
				if um.Id == existing.Id {
					alreadyCloned = true
					break
				}
			}
			if !alreadyCloned {
				umsToBeCloned = append(umsToBeCloned, um)
			}
		}

		for _, umCombo := range umsToBeCloned {
			um, err := memberships.GetUserMembership(umCombo.Id)
			if err != nil {
				return fmt.Errorf("get user membership: %v", err)
			}
			if err := inv.AttachUserMembership(um); err != nil {
				return fmt.Errorf("attach user membership: %v", err)
			}
		}
	}

	if err := o.Begin(); err != nil {
		return fmt.Errorf("begin tx: %v", err)
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

	if err := o.Commit(); err != nil {
		return fmt.Errorf("commit tx: %v", err)
	}

	return
}
