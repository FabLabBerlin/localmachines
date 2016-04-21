package coupons

import (
	"github.com/astaxie/beego/orm"
)

type Coupon struct {
	Id         int64
	LocationId int64
	Code       string
	UserId     int64
	Value      float64
}

type CouponUsage struct {
	Id       int64
	CouponId int64
	Value    float64
	month    int
	year     int
}

func init() {
	orm.RegisterModel(new(Coupon), new(CouponUsage))
}
