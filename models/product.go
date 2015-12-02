package models

import (
	"github.com/astaxie/beego/orm"
)

const (
	PRODUCT_TYPE_CO_WORKING = "co-working"
	PRODUCT_TYPE_SPACE      = "space"
	PRODUCT_TYPE_TUTOR      = "tutor"
)

type Product struct {
	Id            int64
	Type          string `orm:"size(100)"`
	Name          string `orm:"size(100)"`
	Price         float64
	PriceUnit     string `orm:"size(100)"`
	UserId        int64
	MachineSkills string `orm:"size(255)"` // JSON string [1, 3, 6]
	Comments      string `orm:"type(text)"`
	Archived      bool
}

type ProductList struct {
	Data []*Product
}

func (this *Product) TableName() string {
	return "products"
}

func CreateProduct(product *Product) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(product)
}

func GetProduct(productId int64) (product *Product, err error) {
	o := orm.NewOrm()
	product = &Product{Id: productId}
	err = o.Read(product)
	return
}

// Filter out archived products
func GetAllProducts() (products []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Product).TableName()).All(&products)
	return
}

func GetAllProductsOfType(productType string) (products []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Product).TableName()).
		Filter("type", productType).
		Exclude("archived", 1).
		All(&products)
	return
}

func UpdateProduct(product *Product) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(product)
	return
}

func ArchiveProduct(product *Product) (err error) {
	o := orm.NewOrm()
	product.Archived = true
	_, err = o.Update(product)
	return
}

func init() {
	orm.RegisterModel(new(Product))
}
