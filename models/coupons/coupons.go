package coupons

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"time"
)

const CODE_LENGTH = 10

type Coupon struct {
	Id         int64
	LocationId int64
	Code       string
	UserId     int64
	Value      float64
	//TimeEnd    time.Time
}

func GetCoupon(id int64) (c *Coupon, err error) {
	o := orm.NewOrm()
	c = &Coupon{Id: id}
	err = o.Read(c)
	return
}

func GetCouponByCode(locId int64, code string) (*Coupon, error) {
	o := orm.NewOrm()
	var coupon Coupon
	err := o.QueryTable("coupons").
		Filter("location_id", locId).
		Filter("code", code).
		One(&coupon)
	return &coupon, err
}

func GetAllCouponsAt(locId int64) (cs []*Coupon, err error) {
	o := orm.NewOrm()
	c := Coupon{}
	_, err = o.QueryTable(c.TableName()).
		Filter("location_id", locId).
		All(&cs)
	return
}

func GetAllCouponsOf(locId, userId int64) (cs []*Coupon, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("coupons").
		Filter("location_id", locId).
		Filter("user_id", userId).
		All(&cs)
	return
}

func Generate(locId int64, staticCode string, n int, value float64) (cs []*Coupon, err error) {
	cs = make([]*Coupon, 0, n)
	for tries := 0; len(cs) == 0 && tries < 10; tries++ {
		for i := 0; i < n; i++ {
			code := staticCode
			if staticCode == "" {
				if code, err = generateCode(); err != nil {
					return nil, fmt.Errorf("generate code: %v", err)
				}
			}
			c := &Coupon{
				LocationId: locId,
				Code:       code,
				Value:      value,
			}
			cs = append(cs, c)
		}
		if staticCode == "" {
			if err = checkUnique(locId, cs); err != nil {
				cs = make([]*Coupon, 0, n)
			}
		}
	}
	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		return nil, fmt.Errorf("begin: %v", err)
	}
	for _, c := range cs {
		if _, err = o.Insert(c); err != nil {
			o.Rollback()
			return nil, fmt.Errorf("insert: %v", err)
		}
	}
	if err = o.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %v", err)
	}
	return
}

var ErrDuplicateCode = errors.New("duplicate code found in list")

// checkUnique cs.Codes and that none of them is in the db already.
func checkUnique(locId int64, cs []*Coupon) (err error) {
	codes := make(map[string]interface{})
	for _, c := range cs {
		if _, ok := codes[c.Code]; ok {
			return ErrDuplicateCode
		}
		if _, err := GetCouponByCode(locId, c.Code); err == nil {
			return ErrDuplicateCode
		}
	}
	return
}

func generateCode() (code string, err error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, CODE_LENGTH)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b), nil
}

func (c *Coupon) Assign(userId int64) (err error) {
	if c.UserId == userId {
		return nil
	}
	if c.UserId != 0 && c.UserId != userId {
		return errors.New("already assigned to other user")
	}
	c.UserId = userId
	return c.Update()
}

func (c *Coupon) RestValue() (restValue float64, err error) {
	restValue = c.Value
	usages, err := c.Usages()
	if err != nil {
		return 0, fmt.Errorf("usages: %v", err)
	}
	for _, u := range usages {
		restValue -= u.Value
	}
	if restValue < -1 {
		return 0, fmt.Errorf("rest value %v", restValue)
	} else if restValue < 0 {
		restValue = 0
	}
	return
}

func (c *Coupon) TableName() string {
	return "coupons"
}

func (c *Coupon) Update() (err error) {
	_, err = orm.NewOrm().Update(c)
	return
}

func (c *Coupon) Usages() (us []*CouponUsage, err error) {
	_, err = orm.NewOrm().
		QueryTable("coupon_usages").
		Filter("coupon_id", c.Id).
		All(&us)
	return
}

// UseForInvoice so the user gets a rebate to it.
func (c *Coupon) UseForInvoice(invoiceValue float64, month time.Month, year int) (u *CouponUsage, err error) {
	if invoiceValue < 0.01 {
		return
	}
	if c.UserId == 0 {
		return nil, fmt.Errorf("no user assigned: %v", err)
	}
	monthLastDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
	if monthLastDay.After(time.Now()) {
		return nil, fmt.Errorf("coupon cannot be used for %v %v invoice yet",
			month, year)
	}
	usages, err := c.Usages()
	if err != nil {
		return nil, fmt.Errorf("usages: %v", err)
	}
	for _, usage := range usages {
		if usage.Month == int(month) {
			return nil, fmt.Errorf("coupon already used")
		}
	}
	restValue, err := c.RestValue()
	if err != nil {
		return nil, fmt.Errorf("rest value: %v", err)
	}
	if restValue > 0 {
		value := restValue
		if value > invoiceValue {
			value = invoiceValue
		}
		t := time.Now()
		u = &CouponUsage{
			CouponId: c.Id,
			Value:    value,
			Month:    int(t.Month()),
			Year:     t.Year(),
		}
		_, err = orm.NewOrm().Insert(u)
	}
	return
}

type CouponUsage struct {
	Id       int64
	CouponId int64
	Value    float64
	Month    int
	Year     int
}

func (u *CouponUsage) TableName() string {
	return "coupon_usages"
}

func init() {
	orm.RegisterModel(new(Coupon), new(CouponUsage))
}
