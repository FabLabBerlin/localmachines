package models

import (
	"github.com/astaxie/beego/orm"
)

const (
	PRODUCT_TYPE_SPACE = "space"
)

type Product struct {
	Id        int64
	Type      string `orm:"size(100)"`
	Name      string `orm:"size(100)"`
	Price     float64
	PriceUnit string `orm:"size(100)"`
}

func (this *Product) TableName() string {
	return "products"
}

func CreateProduct(product *Product) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(product)
}

func GetAllProducts() (products []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Product).TableName()).All(&products)
	return
}

func UpdateProduct(product *Product) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(product)
	return
}

func init() {
	orm.RegisterModel(new(Product))
}
