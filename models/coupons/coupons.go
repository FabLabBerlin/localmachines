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
}

func (c *Coupon) TableName() string {
	return "coupons"
}

type CouponUsage struct {
	Id       int64
	CouponId int64
	Value    float64
	month    int
	year     int
}

func (cu *CouponUsage) TableName() string {
	return "coupon_usages"
}

func init() {
	orm.RegisterModel(new(Coupon), new(CouponUsage))
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
				Code: code,
				Value: value,
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
		return
	}
	if _, err = o.InsertMulti(n, cs); err != nil {
		o.Rollback()
		return
	}
	err = o.Commit()
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
